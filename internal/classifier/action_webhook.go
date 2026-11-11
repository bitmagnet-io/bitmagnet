package classifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const webhookName = "webhook"

const (
	defaultWebhookTimeout  = 10 * time.Second
	maxWebhookResponseSize = 1 << 20 // 1 MiB
)

type webhookAction struct{}

func (webhookAction) name() string {
	return webhookName
}

type webhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Timeout time.Duration     `json:"timeout,omitempty"`
}

var webhookConfigPayloadSpec = payloadTransformer[webhookConfig, webhookConfig]{
	spec: payloadStruct[webhookConfig]{
		jsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{
					"type":      "string",
					"minLength": 1,
				},
				"headers": map[string]any{
					"type":                 "object",
					"additionalProperties": map[string]any{"type": "string"},
				},
				"timeout": map[string]any{
					"type": "string",
				},
			},
			"required":             []string{"url"},
			"additionalProperties": false,
		},
	},
	transform: func(cfg webhookConfig, _ compilerContext) (webhookConfig, error) {
		u, err := url.Parse(cfg.URL)
		if err != nil {
			return webhookConfig{}, fmt.Errorf("invalid url: %w", err)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return webhookConfig{}, fmt.Errorf("url must use the http or https scheme: %s", cfg.URL)
		}
		if u.Host == "" {
			return webhookConfig{}, fmt.Errorf("url must include a host: %s", cfg.URL)
		}
		return cfg, nil
	},
}

var webhookPayloadSpec = payloadSingleKeyValue[webhookConfig]{
	key:         webhookName,
	valueSpec:   payloadMustSucceed[webhookConfig]{webhookConfigPayloadSpec},
	description: "Send the torrent and the current classification to an external classifier service via HTTP POST",
}

type webhookTorrentFile struct {
	Path      string `json:"path"`
	Size      uint   `json:"size"`
	Extension string `json:"extension,omitempty"`
	FileType  string `json:"fileType,omitempty"`
}

type webhookTorrent struct {
	InfoHash           string               `json:"infoHash"`
	Name               string               `json:"name"`
	BaseName           string               `json:"baseName"`
	Size               uint                 `json:"size"`
	Files              []webhookTorrentFile `json:"files"`
	FileExtensions     []string             `json:"fileExtensions,omitempty"`
	Seeders            *uint                `json:"seeders,omitempty"`
	Leechers           *uint                `json:"leechers,omitempty"`
	HasHint            bool                 `json:"hasHint"`
	HasHintedContentID bool                 `json:"hasHintedContentId"`
}

type webhookContent struct {
	Source        string   `json:"source"`
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Type          string   `json:"type"`
	ReleaseYear   *int     `json:"releaseYear,omitempty"`
	Adult         *bool    `json:"adult,omitempty"`
	OriginalTitle *string  `json:"originalTitle,omitempty"`
	Popularity    *float32 `json:"popularity,omitempty"`
	VoteAverage   *float32 `json:"voteAverage,omitempty"`
	VoteCount     *uint    `json:"voteCount,omitempty"`
}

type webhookClassification struct {
	ContentType        string   `json:"contentType"`
	HasAttachedContent bool     `json:"hasAttachedContent"`
	HasBaseTitle       bool     `json:"hasBaseTitle"`
	BaseTitle          string   `json:"baseTitle,omitempty"`
	Year               *int     `json:"year,omitempty"`
	Languages          []string `json:"languages,omitempty"`
	Episodes           []string `json:"episodes,omitempty"`
	VideoResolution    string   `json:"videoResolution,omitempty"`
	VideoSource        string   `json:"videoSource,omitempty"`
	VideoCodec         string   `json:"videoCodec,omitempty"`
	Video3D            string   `json:"video3d,omitempty"`
	VideoModifier      string   `json:"videoModifier,omitempty"`
	ReleaseGroup       string   `json:"releaseGroup,omitempty"`
	ContentSource      string   `json:"contentSource,omitempty"`
	ContentID          string   `json:"contentId,omitempty"`
	Tags               []string `json:"tags,omitempty"`
}

type webhookRequest struct {
	Torrent        webhookTorrent        `json:"torrent"`
	Classification webhookClassification `json:"classification"`
}

type webhookResponse struct {
	Matched         *bool           `json:"matched,omitempty"`
	Delete          bool            `json:"delete,omitempty"`
	ContentType     *string         `json:"contentType,omitempty"`
	BaseTitle       *string         `json:"baseTitle,omitempty"`
	Year            *int            `json:"year,omitempty"`
	Languages       []string        `json:"languages,omitempty"`
	Episodes        []string        `json:"episodes,omitempty"`
	VideoResolution *string         `json:"videoResolution,omitempty"`
	VideoSource     *string         `json:"videoSource,omitempty"`
	VideoCodec      *string         `json:"videoCodec,omitempty"`
	Video3D         *string         `json:"video3d,omitempty"`
	VideoModifier   *string         `json:"videoModifier,omitempty"`
	ReleaseGroup    *string         `json:"releaseGroup,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
	Content         *webhookContent `json:"content,omitempty"`
}

func (webhookAction) compileAction(ctx compilerContext) (action, error) {
	cfg, err := webhookPayloadSpec.Unmarshal(ctx)
	if err != nil {
		return action{}, ctx.error(err)
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = defaultWebhookTimeout
	}

	client := &http.Client{Timeout: timeout}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return runWebhook(ctx, client, cfg)
		},
	}, nil
}

func (webhookAction) JSONSchema() JSONSchema {
	return webhookPayloadSpec.JSONSchema()
}

func runWebhook(ctx executionContext, client *http.Client, cfg webhookConfig) (classification.Result, error) {
	cl := ctx.result

	reqBody, err := newWebhookRequest(ctx)
	if err != nil {
		return cl, err
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return cl, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL, bytes.NewReader(body))
	if err != nil {
		return cl, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "bitmagnet-classifier/1.0")
	for name, value := range cfg.Headers {
		req.Header.Set(name, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return cl, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxWebhookResponseSize+1))
	if err != nil {
		return cl, err
	}
	if len(respBody) > maxWebhookResponseSize {
		return cl, fmt.Errorf("webhook response exceeds %d bytes", maxWebhookResponseSize)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return cl, fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, truncateBytes(respBody, 256))
	}

	var result webhookResponse
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			return cl, fmt.Errorf("invalid webhook response: %w", err)
		}
	}

	if result.Delete {
		return cl, classification.ErrDeleteTorrent
	}

	if result.Matched != nil && !*result.Matched {
		return cl, classification.ErrUnmatched
	}

	if err := applyWebhookResponse(&cl, result); err != nil {
		return cl, err
	}

	return cl, nil
}

func newWebhookRequest(ctx executionContext) (webhookRequest, error) {
	t := ctx.torrent

	files := make([]webhookTorrentFile, 0, len(t.Files))

	switch t.FilesStatus {
	case model.FilesStatusNoInfo:
		ext := model.FileExtensionFromPath(t.Name)
		if ext.Valid {
			file := webhookTorrentFile{
				Path:      t.Name,
				Extension: ext.String,
			}
			if ft := model.FileTypeFromExtension(ext.String); ft.Valid {
				file.FileType = ft.FileType.String()
			}
			files = append(files, file)
		}
	case model.FilesStatusSingle:
		file := webhookTorrentFile{
			Path: t.Name,
			Size: t.Size,
		}
		if t.Extension.Valid {
			file.Extension = t.Extension.String
		}
		if ft := t.FileType(); ft.Valid {
			file.FileType = ft.FileType.String()
		}
		files = append(files, file)
	case model.FilesStatusMulti, model.FilesStatusOverThreshold:
		for _, f := range t.Files {
			file := webhookTorrentFile{
				Path: f.Path,
				Size: f.Size,
			}
			if f.Extension.Valid {
				file.Extension = f.Extension.String
			}
			if ft := f.FileType(); ft.Valid {
				file.FileType = ft.FileType.String()
			}
			files = append(files, file)
		}
	}

	req := webhookRequest{
		Torrent: webhookTorrent{
			InfoHash:           t.InfoHash.String(),
			Name:               t.Name,
			BaseName:           t.BaseName(),
			Size:               t.Size,
			Files:              files,
			FileExtensions:     t.FileExtensions(),
			HasHint:            !t.Hint.IsNil(),
			HasHintedContentID: !t.Hint.IsNil() && t.Hint.ContentSource.Valid,
		},
		Classification: newWebhookClassification(ctx.result),
	}

	if seeders := t.Seeders(); seeders.Valid {
		req.Torrent.Seeders = &seeders.Uint
	}
	if leechers := t.Leechers(); leechers.Valid {
		req.Torrent.Leechers = &leechers.Uint
	}

	return req, nil
}

func newWebhookClassification(cl classification.Result) webhookClassification {
	contentType := "unknown"
	if cl.ContentType.Valid {
		contentType = cl.ContentType.ContentType.String()
	}

	c := webhookClassification{
		ContentType:        contentType,
		HasAttachedContent: cl.Content != nil,
		HasBaseTitle:       cl.BaseTitle.Valid,
		ReleaseGroup:       cl.ReleaseGroup.String,
	}

	if cl.BaseTitle.Valid {
		c.BaseTitle = cl.BaseTitle.String
	}
	if cl.Content != nil {
		c.ContentSource = cl.Content.Source
		c.ContentID = cl.Content.ID
	}
	if cl.Date.Year > 0 {
		year := int(cl.Date.Year)
		c.Year = &year
	}
	for lang := range cl.Languages {
		c.Languages = append(c.Languages, lang.String())
	}
	for _, s := range cl.Episodes.SeasonEntries() {
		c.Episodes = append(c.Episodes, s.String())
	}
	if cl.VideoResolution.Valid {
		c.VideoResolution = cl.VideoResolution.VideoResolution.String()
	}
	if cl.VideoSource.Valid {
		c.VideoSource = cl.VideoSource.VideoSource.String()
	}
	if cl.VideoCodec.Valid {
		c.VideoCodec = cl.VideoCodec.VideoCodec.String()
	}
	if cl.Video3D.Valid {
		c.Video3D = cl.Video3D.Video3D.String()
	}
	if cl.VideoModifier.Valid {
		c.VideoModifier = cl.VideoModifier.VideoModifier.String()
	}
	for tag := range cl.Tags {
		c.Tags = append(c.Tags, tag)
	}

	return c
}

func applyWebhookResponse(cl *classification.Result, resp webhookResponse) error {
	if resp.ContentType != nil {
		if *resp.ContentType == "unknown" {
			cl.ContentType = model.NullContentType{}
		} else {
			contentType, err := model.ParseContentType(*resp.ContentType)
			if err != nil {
				return fmt.Errorf("invalid content type '%s': %w", *resp.ContentType, err)
			}
			cl.ContentType = model.NewNullContentType(contentType)
		}
	}

	if resp.BaseTitle != nil {
		cl.BaseTitle = model.NewNullString(*resp.BaseTitle)
	}

	if resp.Year != nil {
		cl.Date = model.Date{Year: model.Year(*resp.Year)}
	}

	if resp.Languages != nil {
		languages := make(model.Languages)
		for _, name := range resp.Languages {
			lang := model.ParseLanguage(name)
			if !lang.Valid {
				return fmt.Errorf("invalid language '%s'", name)
			}
			languages[lang.Language] = struct{}{}
		}
		cl.Languages = languages
	}

	if resp.Episodes != nil {
		episodes := model.Episodes{}
		for _, s := range resp.Episodes {
			parsed := model.ParseEpisodes(s)
			if parsed == nil {
				return fmt.Errorf("invalid episodes '%s'", s)
			}
			for season, eps := range parsed {
				if episodes[season] == nil {
					episodes[season] = make(map[int]struct{})
				}
				for ep := range eps {
					episodes[season][ep] = struct{}{}
				}
			}
		}
		cl.Episodes = episodes
	}

	if resp.VideoResolution != nil {
		resolution, err := model.ParseVideoResolution(*resp.VideoResolution)
		if err != nil {
			return fmt.Errorf("invalid video resolution '%s': %w", *resp.VideoResolution, err)
		}
		cl.VideoResolution = model.NewNullVideoResolution(resolution)
	}

	if resp.VideoSource != nil {
		source, err := model.ParseVideoSource(*resp.VideoSource)
		if err != nil {
			return fmt.Errorf("invalid video source '%s': %w", *resp.VideoSource, err)
		}
		cl.VideoSource = model.NewNullVideoSource(source)
	}

	if resp.VideoCodec != nil {
		codec, err := model.ParseVideoCodec(*resp.VideoCodec)
		if err != nil {
			return fmt.Errorf("invalid video codec '%s': %w", *resp.VideoCodec, err)
		}
		cl.VideoCodec = model.NewNullVideoCodec(codec)
	}

	if resp.Video3D != nil {
		threeD, err := model.ParseVideo3D(*resp.Video3D)
		if err != nil {
			return fmt.Errorf("invalid video 3d '%s': %w", *resp.Video3D, err)
		}
		cl.Video3D = model.NewNullVideo3D(threeD)
	}

	if resp.VideoModifier != nil {
		modifier, err := model.ParseVideoModifier(*resp.VideoModifier)
		if err != nil {
			return fmt.Errorf("invalid video modifier '%s': %w", *resp.VideoModifier, err)
		}
		cl.VideoModifier = model.NewNullVideoModifier(modifier)
	}

	if resp.ReleaseGroup != nil {
		cl.ReleaseGroup = model.NewNullString(*resp.ReleaseGroup)
	}

	if resp.Tags != nil {
		tags := make(map[string]struct{}, len(resp.Tags))
		for _, tag := range resp.Tags {
			tags[tag] = struct{}{}
		}
		cl.Tags = tags
	}

	if resp.Content != nil {
		content, err := newWebhookContent(*resp.Content)
		if err != nil {
			return err
		}
		cl.AttachContent(&content)
	}

	return nil
}

func newWebhookContent(c webhookContent) (model.Content, error) {
	if c.Source == "" || c.ID == "" {
		return model.Content{}, fmt.Errorf("content.source and content.id are required")
	}

	contentType, err := model.ParseContentType(c.Type)
	if err != nil {
		return model.Content{}, fmt.Errorf("invalid content type '%s': %w", c.Type, err)
	}

	content := model.Content{
		Source: c.Source,
		ID:     c.ID,
		Title:  c.Title,
		Type:   contentType,
	}
	if c.ReleaseYear != nil {
		content.ReleaseYear = model.Year(*c.ReleaseYear)
	}
	if c.Adult != nil {
		content.Adult = model.NewNullBool(*c.Adult)
	}
	if c.OriginalTitle != nil {
		content.OriginalTitle = model.NewNullString(*c.OriginalTitle)
	}
	if c.Popularity != nil {
		content.Popularity = model.NewNullFloat32(*c.Popularity)
	}
	if c.VoteAverage != nil {
		content.VoteAverage = model.NewNullFloat32(*c.VoteAverage)
	}
	if c.VoteCount != nil {
		content.VoteCount = model.NewNullUint(*c.VoteCount)
	}

	return content, nil
}

func truncateBytes(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
