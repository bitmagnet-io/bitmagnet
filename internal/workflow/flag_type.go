package workflow

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -f flag_type.go

// FlagType represents the type of a flag
// ENUM(bool, string, int, string_list, content_type_list)
type FlagType string
