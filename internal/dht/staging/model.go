package staging

import (
	"errors"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"unicode/utf8"
)

func createTorrentModel(
	infoHash krpc.ID,
	info metainfo.Info,
	peers ResponseScrape,
) (model.Torrent, error) {
	name := info.BestName()
	checkUtf8Strings := []string{name}
	private := false
	if info.Private != nil {
		private = *info.Private
	}
	var files []model.TorrentFile
	for i, file := range info.Files {
		displayPath := file.DisplayPath(&info)
		files = append(files, model.TorrentFile{
			Index: uint32(i),
			Path:  displayPath,
			Size:  uint64(file.Length),
		})
		checkUtf8Strings = append(checkUtf8Strings, displayPath)
	}
	for _, str := range checkUtf8Strings {
		if !utf8.ValidString(str) {
			return model.Torrent{}, errors.New("invalid utf8 string")
		}
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
		InfoHash:    model.Hash20(infoHash),
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
	infoHash krpc.ID,
	peers ResponseScrape,
) (model.TorrentsTorrentSource, error) {
	seeders := model.NullUint{}
	leechers := model.NullUint{}
	bfsdSize := uint(peers.Bfsd.ApproximatedSize())
	bfpeSize := uint(peers.Bfpe.ApproximatedSize())
	if peers.Scraped || bfsdSize > 0 {
		seeders.Valid = true
		seeders.Uint = bfsdSize
	}
	if peers.Scraped || bfpeSize > 0 {
		leechers.Valid = true
		leechers.Uint = bfpeSize
	}
	// todo add discovered peers to bloom
	bfsdBytes, bfsdErr := peers.Bfsd.MarshalBinary()
	if bfsdErr != nil {
		return model.TorrentsTorrentSource{}, bfsdErr
	}
	bfpeBytes, bfpeErr := peers.Bfpe.MarshalBinary()
	if bfpeErr != nil {
		return model.TorrentsTorrentSource{}, bfpeErr
	}
	return model.TorrentsTorrentSource{
		Source:   "dht",
		InfoHash: model.Hash20(infoHash),
		Bfsd:     bfsdBytes,
		Bfpe:     bfpeBytes,
		Leechers: seeders,
		Seeders:  leechers,
	}, nil
}
