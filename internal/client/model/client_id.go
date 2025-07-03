package model

//go:generate go run github.com/abice/go-enum --marshal --names --nocomments --sql --sqlnullstr --values -f client_id.go

// ID represents client that has an implemented interface
// ENUM(QBittorrent, Transmission, Ntfy)
type ID string
