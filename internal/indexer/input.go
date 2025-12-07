package indexer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

func InputContent(cs ...model.Content) Input {
	return func(p *payload) {
		for _, c := range cs {
			p.content.Set(c.Ref(), c)
		}
	}
}

func InputTorrentContent(tcs ...model.TorrentContent) Input {
	return func(p *payload) {
		for _, tc := range tcs {
			p.torrentContent.Set(tc.Ref(), tc)

			if tc.Content.ID != "" {
				InputContent(tc.Content)(p)
			}
		}
	}
}

func InputDeleteTorrentContent(refs ...model.TorrentContentRef) Input {
	return func(p *payload) {
		p.deleteTorrentContent.SetKeys(refs...)
	}
}

func InputDeleteInfoHashes(infoHashes ...protocol.ID) Input {
	return func(p *payload) {
		p.deleteInfoHash.SetKeys(infoHashes...)
	}
}

type Inputs []Input

func (i Inputs) Input() Input {
	return func(p *payload) {
		for _, i := range i {
			i(p)
		}
	}
}

type payload struct {
	content              maps.InsertMap[model.ContentRef, model.Content]
	torrentContent       maps.InsertMap[model.TorrentContentRef, model.TorrentContent]
	deleteTorrentContent maps.InsertMap[model.TorrentContentRef, struct{}]
	deleteInfoHash       maps.InsertMap[protocol.ID, struct{}]
}

func newPayload() payload {
	return payload{
		content:              maps.NewInsertMap[model.ContentRef, model.Content](),
		torrentContent:       maps.NewInsertMap[model.TorrentContentRef, model.TorrentContent](),
		deleteTorrentContent: maps.NewInsertMap[model.TorrentContentRef, struct{}](),
		deleteInfoHash:       maps.NewInsertMap[protocol.ID, struct{}](),
	}
}
