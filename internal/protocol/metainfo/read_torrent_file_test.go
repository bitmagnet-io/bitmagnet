package metainfo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTorrentFile(t *testing.T) {
	input, readErr := os.ReadFile("examples/ubuntu-23.04-desktop-amd64.iso.torrent")
	if readErr != nil {
		t.Fatalf("error reading torrent file: %s", readErr)
	}

	torrentFile, parseErr := ReadTorrentFileBytes(input)
	if parseErr != nil {
		t.Fatalf("error parsing torrent file: %s", parseErr)
	}

	assert.NotEmpty(t, torrentFile.Info.Pieces)
	torrentFile.Info.Pieces = nil
	assert.Equal(t, TorrentFile{
		Info: Info{
			Name:        "ubuntu-23.04-desktop-amd64.iso",
			PieceLength: 262144,
			Length:      4932407296,
		},
		Announce: "https://torrent.ubuntu.com/announce",
		AnnounceList: [][]string{
			{"https://torrent.ubuntu.com/announce"},
			{"https://ipv6.torrent.ubuntu.com/announce"},
		},
		CreationDate: 1681992794,
		Comment:      "Ubuntu CD releases.ubuntu.com",
		CreatedBy:    "mktorrent 1.1",
	}, torrentFile)
}
