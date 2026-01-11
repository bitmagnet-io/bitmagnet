package search

//revive:disable:line-length-limit
//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -f order.go -f order_torrent_content.go -f order_torrent_files.go -f order_queue_jobs.go -f facet.go -f result_type.go
