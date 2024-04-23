package model

import "strings"

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -t enums.gql.tmpl -f content_type.go -f facet_logic.go -f file_type.go -f files_status.go -f queue_job_status.go -f video_3d.go -f video_codec.go -f video_modifier.go -f video_resolution.go -f video_source.go

func removeEnumPrefixes(names ...string) []string {
	result := make([]string, len(names))
	for i, name := range names {
		result[i] = name[1:]
	}
	return result
}

func namesToLower(names ...string) []string {
	result := make([]string, len(names))
	for i, name := range names {
		result[i] = strings.ToLower(name)
	}
	return result
}
