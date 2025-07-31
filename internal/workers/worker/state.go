package worker

//revive:disable:line-length-limit

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sqlnullstr --values -t ../../gql/enums.gql.tmpl -f state.go

// State represents the state of a Worker
// ENUM(idle, startup, running, shutdown, error)
type State string
