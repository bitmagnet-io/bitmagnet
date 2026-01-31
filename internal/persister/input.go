package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Input func(*payload)

type Inputs []Input

func (i Inputs) Input() Input {
	return func(p *payload) {
		for _, input := range i {
			input(p)
		}
	}
}

func InputTorrentSources(torrentSources ...model.TorrentSource) Input {
	return func(p *payload) {
		p.torrentSources.SetEntries(
			slice.Map(
				torrentSources,
				func(ts model.TorrentSource) maps.MapEntry[string, model.TorrentSource] {
					return maps.MapEntry[string, model.TorrentSource]{
						Key:   ts.Key,
						Value: ts,
					}
				},
			)...)
	}
}

func torrentsTorrentSourcesEntries(
	torrentsTorrentSources ...model.TorrentsTorrentSource,
) []maps.MapEntry[hashWithID, model.TorrentsTorrentSource] {
	return slice.Map(
		torrentsTorrentSources,
		func(ts model.TorrentsTorrentSource) maps.MapEntry[hashWithID, model.TorrentsTorrentSource] {
			return maps.MapEntry[hashWithID, model.TorrentsTorrentSource]{
				Key: hashWithID{
					hash: ts.InfoHash,
					id:   ts.Source,
				},
				Value: ts,
			}
		},
	)
}

func InputTorrentsTorrentSources(torrentsTorrentSources ...model.TorrentsTorrentSource) Input {
	return func(p *payload) {
		p.torrentsTorrentSources.SetEntries(
			torrentsTorrentSourcesEntries(torrentsTorrentSources...)...,
		)
	}
}

func torrentPiecesEntries(torrentPieces ...model.TorrentPieces) []maps.MapEntry[protocol.ID, model.TorrentPieces] {
	return slice.Map(torrentPieces, func(tp model.TorrentPieces) maps.MapEntry[protocol.ID, model.TorrentPieces] {
		return maps.MapEntry[protocol.ID, model.TorrentPieces]{
			Key:   tp.InfoHash,
			Value: tp,
		}
	})
}

func InputTorrentPieces(torrentPieces ...model.TorrentPieces) Input {
	return func(p *payload) {
		p.torrentPieces.SetEntries(torrentPiecesEntries(torrentPieces...)...)
	}
}

func torrentFilesEntries(torrentFiles ...model.TorrentFile) []maps.MapEntry[hashWithID, model.TorrentFile] {
	return slice.Map(
		torrentFiles,
		func(tf model.TorrentFile) maps.MapEntry[hashWithID, model.TorrentFile] {
			return maps.MapEntry[hashWithID, model.TorrentFile]{
				Key: hashWithID{
					hash: tf.InfoHash,
					id:   tf.Path,
				},
				Value: tf,
			}
		})
}

func InputTorrentFiles(torrentFiles ...model.TorrentFile) Input {
	return func(p *payload) {
		p.torrentFiles.SetEntries(torrentFilesEntries(torrentFiles...)...)
	}
}

func InputTorrents(torrents ...model.Torrent) Input {
	return func(p *payload) {
		p.torrents.SetEntries(
			slice.Map(torrents, func(t model.Torrent) maps.MapEntry[protocol.ID, model.Torrent] {
				return maps.MapEntry[protocol.ID, model.Torrent]{
					Key:   t.InfoHash,
					Value: t,
				}
			})...)
	}
}

func InputContent(content ...model.Content) Input {
	return func(p *payload) {
		p.content.SetEntries(slice.Map(
			content,
			func(c model.Content) maps.MapEntry[model.ContentRef, model.Content] {
				return maps.MapEntry[model.ContentRef, model.Content]{
					Key:   c.Ref(),
					Value: c,
				}
			},
		)...)
	}
}

func InputTorrentContent(torrentContents ...model.TorrentContent) Input {
	return func(p *payload) {
		p.torrentContents.SetEntries(slice.Map(
			torrentContents,
			func(t model.TorrentContent) maps.MapEntry[model.TorrentContentRef, model.TorrentContent] {
				return maps.MapEntry[model.TorrentContentRef, model.TorrentContent]{
					Key:   t.Ref(),
					Value: t,
				}
			},
		)...)
	}
}

func InputDeleteTorrentContent(torrentContentRefs ...model.TorrentContentRef) Input {
	return func(p *payload) {
		p.deleteTorrentContent.SetEntries(
			slice.Map(
				torrentContentRefs,
				func(ref model.TorrentContentRef) maps.MapEntry[model.TorrentContentRef, struct{}] {
					return maps.MapEntry[model.TorrentContentRef, struct{}]{
						Key: ref,
					}
				},
			)...)
	}
}

func InputDeleteInfoHashes(infoHashes ...protocol.ID) Input {
	return func(p *payload) {
		p.deleteInfoHashes.SetEntries(
			slice.Map(infoHashes, func(id protocol.ID) maps.MapEntry[protocol.ID, struct{}] {
				return maps.MapEntry[protocol.ID, struct{}]{
					Key: id,
				}
			})...)
	}
}

func InputTorrentTags(infoHash protocol.ID, tags ...string) Input {
	return func(p *payload) {
		p.torrentTags.SetEntries(slice.Map(
			tags,
			func(tagName string) maps.MapEntry[hashWithID, bool] {
				return maps.MapEntry[hashWithID, bool]{
					Key: hashWithID{
						hash: infoHash,
						id:   tagName,
					},
					Value: true,
				}
			})...)
	}
}

func InputDeleteTorrentTags(infoHash protocol.ID, tags ...string) Input {
	return func(p *payload) {
		p.torrentTags.SetEntries(slice.Map(
			tags,
			func(tagName string) maps.MapEntry[hashWithID, bool] {
				return maps.MapEntry[hashWithID, bool]{
					Key: hashWithID{
						hash: infoHash,
						id:   tagName,
					},
					Value: false,
				}
			})...)
	}
}

func InputQueueJobs(queueJobs ...model.QueueJob) Input {
	return func(p *payload) {
		p.queueJobs.SetEntries(
			slice.Map(queueJobs, func(j model.QueueJob) maps.MapEntry[string, model.QueueJob] {
				return maps.MapEntry[string, model.QueueJob]{
					Key:   j.Queue + ":" + j.Fingerprint,
					Value: j,
				}
			})...)
	}
}

func InputFlush(p *payload) {
	p.shouldFlush = true
}
