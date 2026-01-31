package plugin

const KeyActivation = "activation"

// revive:disable:line-length-limit
//go:generate go run github.com/abice/go-enum --marshal --names --nocomments --sql --sqlnullstr --values -t ../../internal/gql/enums.gql.tmpl -f activation.go

// Activation represents the activation mode for a plugin.
/* ENUM(enabled, disabled, auto, always) */
type Activation string
