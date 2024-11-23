package classifier

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"slices"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type pluginContentList []struct {
	SeriesId *interface{} `json:"seriesId,omitempty"`
	TmdbId   *int         `json:"tmdbId,omitempty"`
	Id       *interface{} `json:"Id,omitempty"`
}

func (l pluginContentList) TmdbIds(p pluginSource) []string {
	var list []string
	for _, item := range l {
		if item.TmdbId != nil && *item.TmdbId > 0 {
			list = append(list, strconv.Itoa(*item.TmdbId))
		} else if item.TmdbId == nil && (item.SeriesId != nil || item.Id != nil) {
			// there are cases of episodes and series
			if item.SeriesId == nil {
				item.SeriesId = item.Id
			}
			u, _ := url.Parse(p.Url)
			var tmdbid struct {
				TmdbId      int `json:"tmdbId"`
				ProviderIds struct {
					Tmdb string `json:"tmdb"`
				} `json:"ProviderIds"`
			}
			req := resty.New().R().SetQueryParam("apikey", *p.ApiKey).SetResult(&tmdbid)
			switch (*item.SeriesId).(type) {
			case string:
				req.Get(fmt.Sprintf("%s://%s/Items/%v?%s", u.Scheme, u.Host, *item.SeriesId, u.RawQuery))
				list = append(list, tmdbid.ProviderIds.Tmdb)
			default:
				req.Get(fmt.Sprintf("%s://%s%s/series/%v", u.Scheme, u.Host, path.Dir(u.Path), *item.SeriesId))
				if tmdbid.TmdbId > 0 {
					list = append(list, strconv.Itoa(tmdbid.TmdbId))
				}
			}
		}
	}

	slices.Sort(list)
	return slices.Compact(list)
}

type pluginContentStruct struct {
	Page    int `json:"page"`
	Results []struct {
		Id int `json:"id"`
	} `json:"results"`
	TotalPages int `json:"total_pages"`
}

func (l pluginContentStruct) TmdbIds() []string {
	var list []string
	for _, item := range l.Results {
		list = append(list, strconv.Itoa(item.Id))
	}

	return list
}

// a string_list has to be []any not []string
func (p pluginSource) any(data []string) []any {
	anydata := make([]any, len(data))
	for i, v := range data {
		anydata[i] = v
	}
	return anydata
}

func (p pluginSource) source() ([]string, error) {
	var r pluginContentList
	var s pluginContentStruct
	req := resty.New().R()
	if p.ApiKey != nil {
		u, _ := url.Parse(p.Url)
		parmName := "apikey"
		switch u.Host {
		case "api.themoviedb.org":
			parmName = "api_key"
		}
		req = req.SetQueryParam(parmName, *p.ApiKey)
	}
	if p.Start != nil {
		req = req.SetQueryParam(*p.Start, time.Now().Add(time.Duration(*p.Days)*24*time.Hour*-1).Format(time.DateOnly))
	}
	if p.End != nil {
		req = req.SetQueryParam(*p.End, time.Now().Add(time.Duration(*p.Days+1)*24*time.Hour).Format(time.DateOnly))
	}
	resp, err := req.Get(p.Url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("%v %v\n%s\n%s\n", resp.StatusCode(), resp.Request.URL, resp.Status(), resp.Body())
	}
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("%v %v %v", resp.StatusCode(), resp.Request.URL, resp.Status())

	}

	// try type of response from sonarr / radarr / jellyfin
	err = json.Unmarshal(resp.Body(), &r)
	if err == nil {
		return r.TmdbIds(p), nil

	}

	// try type of response from tmdb
	err = json.Unmarshal(resp.Body(), &s)
	if err == nil {
		list := s.TmdbIds()
		for s.Page < s.TotalPages {
			_, err = req.SetQueryParam("page", strconv.Itoa(s.Page+1)).SetResult(&s).Get(p.Url)
			if err != nil {
				return list, err
			}
			list = append(list, s.TmdbIds()...)

		}
		return list, nil
	}

	return nil, err
}
