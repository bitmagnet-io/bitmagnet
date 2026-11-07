// Package matchlab provides the only transport for MatchLab's public API.
// It deliberately keeps title identification separate from TMDB enrichment.
package matchlab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	SIDStart  *string `json:"sid_start"`
	SIDEnd    *string `json:"sid_end"`
	ASIDStart *string `json:"asid_start"`
	ASIDEnd   *string `json:"asid_end"`
}

// UnmarshalJSON enforces the public /match contract. A partial or empty
// response must never become a legitimate no-match and clear prior matches.
func (r *Response) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) != 4 {
		return fmt.Errorf("MatchLab /match requires exactly four identity fields, got %d", len(fields))
	}
	for _, name := range []string{"sid_start", "sid_end", "asid_start", "asid_end"} {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("MatchLab /match missing identity field %q", name)
		}
	}
	type plain Response
	var parsed plain
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	*r = Response(parsed)
	return nil
}

type Parsed struct {
	Title string `json:"title"`
	Year  *int   `json:"year"`
}

// Torrent title parsers may encode year as a JSON number or a string.
// Do not silently discard otherwise valid title metadata on type mismatch.
func (p *Parsed) UnmarshalJSON(data []byte) error {
	var value struct {
		Title string          `json:"title"`
		Year  json.RawMessage `json:"year"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	p.Title = value.Title
	p.Year = nil
	if len(value.Year) == 0 || string(value.Year) == "null" {
		return nil
	}
	var number int
	if err := json.Unmarshal(value.Year, &number); err == nil {
		p.Year = &number
		return nil
	}
	var raw string
	if err := json.Unmarshal(value.Year, &raw); err != nil {
		return fmt.Errorf("invalid MatchLab year: %w", err)
	}
	if number, err := strconv.Atoi(raw); err == nil {
		p.Year = &number
	}
	return nil
}

type Client interface {
	Match(ctx context.Context, title string) (Response, bool, error)
	Parse(ctx context.Context, title string) (Parsed, error)
	Health(ctx context.Context) error
}

var ErrInvalidResponse = errors.New("invalid MatchLab response")

type httpClient struct {
	base  url.URL
	http  *http.Client
	slots chan struct{}
}

func NewClient(cfg Config) (Client, error) {
	if !cfg.Enabled() {
		//nolint:nilnil // An absent URL intentionally disables the optional client.
		return nil, nil
	}
	raw := strings.TrimSpace(cfg.URL)
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid MATCHLAB_URL: %w", err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" ||
		u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
		return nil, errors.New(
			"MATCHLAB_URL must be an HTTP(S) base URL without credentials, query or fragment",
		)
	}
	if cfg.TimeoutSeconds <= 0 || cfg.MaxConcurrent <= 0 {
		return nil, errors.New("MatchLab timeout and concurrency must be positive")
	}
	u.Path = strings.TrimRight(u.Path, "/")
	return &httpClient{
		base: *u,
		http: &http.Client{
			Timeout:       time.Duration(cfg.TimeoutSeconds) * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		slots: make(chan struct{}, cfg.MaxConcurrent),
	}, nil
}

func (c *httpClient) request(ctx context.Context, path string, title string, out any) error {
	select {
	case c.slots <- struct{}{}:
		defer func() { <-c.slots }()
	case <-ctx.Done():
		return ctx.Err()
	}
	u := c.base
	u.Path = strings.TrimSuffix(u.Path, "/") + path
	if title != "" {
		query := u.Query()
		query.Set("title", title)
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("MatchLab %s request: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MatchLab %s HTTP %d", path, resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 32*1024))
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidResponse, err)
	}
	// An otherwise valid identity followed by another JSON value must not
	// be silently accepted as an authoritative MatchLab response.
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("%w: multiple JSON values", ErrInvalidResponse)
		}
		return fmt.Errorf("%w: invalid trailing JSON: %w", ErrInvalidResponse, err)
	}
	return nil
}

func (c *httpClient) Match(ctx context.Context, title string) (Response, bool, error) {
	var r Response
	if err := c.request(ctx, "/match", title, &r); err != nil {
		return Response{}, false, err
	}
	if err := validateResponse(r); err != nil {
		return Response{}, false, err
	}
	return r, r.SIDStart != nil || r.ASIDStart != nil || r.ASIDEnd != nil, nil
}

func (c *httpClient) Parse(ctx context.Context, title string) (Parsed, error) {
	var result Parsed
	err := c.request(ctx, "/parse", title, &result)
	return result, err
}

func (c *httpClient) Health(ctx context.Context) error {
	return c.request(ctx, "/health/ready", "", nil)
}

func ValidIMDb(id string) bool {
	if len(id) < 9 || len(id) > 12 || !strings.HasPrefix(id, "tt") {
		return false
	}
	for _, b := range id[2:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

func imdbSID(value *string) (string, int, bool) {
	if value == nil {
		return "", 0, false
	}
	parts := strings.Split(*value, ":")
	if len(parts) < 1 || len(parts) > 3 || !ValidIMDb(parts[0]) {
		return "", 0, false
	}
	for _, part := range parts[1:] {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 || n > 99999 {
			return "", 0, false
		}
	}
	return parts[0], len(parts), true
}

func validASID(v *string) bool {
	if v == nil {
		return true
	}
	parts := strings.Split(*v, ":")
	if len(parts) != 2 {
		return false
	}
	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return false
		}
	}
	return true
}

func validateResponse(r Response) error {
	if !validASID(r.ASIDStart) || !validASID(r.ASIDEnd) {
		return fmt.Errorf("%w: malformed AniDB coordinates", ErrInvalidResponse)
	}
	if r.SIDStart == nil && r.SIDEnd == nil {
		return nil
	}
	startID, startParts, startOK := imdbSID(r.SIDStart)
	endID, endParts, endOK := imdbSID(r.SIDEnd)
	if !startOK || !endOK || startID != endID || startParts != endParts {
		return fmt.Errorf("%w: inconsistent IMDb coordinates", ErrInvalidResponse)
	}
	return nil
}

func IMDbID(r Response) (string, bool) {
	id, _, ok := imdbSID(r.SIDStart)
	return id, ok
}

func (r Response) IsSeries() bool {
	if r.ASIDStart != nil || r.ASIDEnd != nil {
		return true
	}
	_, parts, ok := imdbSID(r.SIDStart)
	return ok && parts > 1
}
