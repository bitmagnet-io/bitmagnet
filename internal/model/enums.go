package model

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -t enums.gql.tmpl -f content_type.go -f facet_logic.go -f file_type.go -f files_status.go -f video_3d.go -f video_codec.go -f video_modifier.go -f video_resolution.go -f video_source.go

func removeEnumPrefixes(names ...string) []string {
	var result []string
	for _, name := range names {
		result = append(result, name[1:])
	}
	return result
}
