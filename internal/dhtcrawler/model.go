package dhtcrawler

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
)

func createTorrentModel(
	infoHash protocol.ID,
	info metainfo.Info,
	peers stagingResponseScrape,
) (model.Torrent, error) {
	name := info.BestName()
	private := false
	if info.Private != nil {
		private = *info.Private
	}
	var files []model.TorrentFile
	for i, file := range info.Files {
		files = append(files, model.TorrentFile{
			Index: uint32(i),
			Path:  file.DisplayPath(&info),
			Size:  uint64(file.Length),
		})
	}
	source, sourceErr := createTorrentSourceModel(infoHash, peers)
	if sourceErr != nil {
		return model.Torrent{}, sourceErr
	}
	filesStatus := model.FilesStatusSingle
	if len(files) > 0 {
		filesStatus = model.FilesStatusMulti
	}
	return model.Torrent{
		InfoHash:    infoHash,
		Name:        name,
		Size:        uint64(info.TotalLength()),
		Private:     private,
		PieceLength: model.NewNullUint64(uint64(info.PieceLength)),
		Pieces:      info.Pieces,
		Files:       files,
		FilesStatus: filesStatus,
		Sources: []model.TorrentsTorrentSource{
			source,
		},
	}, nil
}

func createTorrentSourceModel(
	infoHash protocol.ID,
	peers stagingResponseScrape,
) (model.TorrentsTorrentSource, error) {
	seeders := model.NullUint{}
	leechers := model.NullUint{}
	bfsdSize := uint(peers.bfsd.ApproximatedSize())
	bfpeSize := uint(peers.bfpe.ApproximatedSize())
	if peers.scraped || bfsdSize > 0 {
		seeders.Valid = true
		seeders.Uint = bfsdSize
	}
	if peers.scraped || bfpeSize > 0 {
		leechers.Valid = true
		leechers.Uint = bfpeSize
	}
	// todo add discovered peers to bloom
	bfsdBytes, bfsdErr := peers.bfsd.MarshalBinary()
	if bfsdErr != nil {
		return model.TorrentsTorrentSource{}, bfsdErr
	}
	bfpeBytes, bfpeErr := peers.bfpe.MarshalBinary()
	if bfpeErr != nil {
		return model.TorrentsTorrentSource{}, bfpeErr
	}
	return model.TorrentsTorrentSource{
		Source:   "dht",
		InfoHash: infoHash,
		Bfsd:     bfsdBytes,
		Bfpe:     bfpeBytes,
		Leechers: seeders,
		Seeders:  leechers,
	}, nil
}
