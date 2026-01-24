module github.com/bitmagnet-io/plugin-qbittorrent

go 1.25.1

require (
	github.com/bitmagnet-io/bitmagnet v0.0.0
	github.com/nicksnyder/go-i18n/v2 v2.6.0
	golang.org/x/text v0.28.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/tetratelabs/wazero v1.9.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/bitmagnet-io/bitmagnet => ../..
