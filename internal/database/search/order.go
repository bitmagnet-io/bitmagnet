package search

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -f order.go -f order_torrent_content.go -f order_torrent_files.go -f order_queue_jobs.go

// OrderDirection represents sort order directions
// ENUM(Ascending, Descending)
type OrderDirection string
