package resolvers

import (
	"context"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/stretchr/testify/require"
)

type fakeSearch struct {
	torrentFilesCalls  int
	torrentFilesResult search.TorrentFilesResult
	torrentFilesErr    error
}

func (f *fakeSearch) Content(_ context.Context, _ ...query.Option) (search.ContentResult, error) {
	return search.ContentResult{}, nil
}

func (f *fakeSearch) QueueJobs(_ context.Context, _ ...query.Option) (search.QueueJobResult, error) {
	return search.QueueJobResult{}, nil
}

func (f *fakeSearch) Torrents(_ context.Context, _ ...query.Option) (search.TorrentsResult, error) {
	return search.TorrentsResult{}, nil
}

func (f *fakeSearch) TorrentsWithMissingInfoHashes(
	_ context.Context,
	_ []protocol.ID,
	_ ...query.Option,
) (search.TorrentsWithMissingInfoHashesResult, error) {
	return search.TorrentsWithMissingInfoHashesResult{}, nil
}

func (f *fakeSearch) TorrentSuggestTags(
	_ context.Context,
	_ search.SuggestTagsQuery,
	_ ...query.Option,
) (search.TorrentSuggestTagsResult, error) {
	return search.TorrentSuggestTagsResult{}, nil
}

func (f *fakeSearch) TorrentContent(_ context.Context, _ ...query.Option) (search.TorrentContentResult, error) {
	return search.TorrentContentResult{}, nil
}

func (f *fakeSearch) TorrentFiles(_ context.Context, _ ...query.Option) (search.TorrentFilesResult, error) {
	f.torrentFilesCalls++
	return f.torrentFilesResult, f.torrentFilesErr
}

func TestTorrentResolverFilesReturnsPreloadedFiles(t *testing.T) {
	fake := &fakeSearch{}
	resolver := torrentResolver{&Resolver{Search: fake}}
	expected := []model.TorrentFile{
		{
			InfoHash: protocol.MustParseID("1111111111111111111111111111111111111111"),
			Index:    0,
			Path:     "movie.mkv",
			Size:     1024,
		},
	}

	got, err := resolver.Files(context.Background(), &model.Torrent{Files: expected})
	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Zero(t, fake.torrentFilesCalls)
}

func TestTorrentResolverFilesReturnsNilWithoutFilesInfo(t *testing.T) {
	fake := &fakeSearch{}
	resolver := torrentResolver{&Resolver{Search: fake}}

	got, err := resolver.Files(context.Background(), &model.Torrent{})
	require.NoError(t, err)
	require.Nil(t, got)
	require.Zero(t, fake.torrentFilesCalls)
}

func TestTorrentResolverFilesLoadsMissingFilesFromSearch(t *testing.T) {
	expected := []model.TorrentFile{
		{
			InfoHash: protocol.MustParseID("2222222222222222222222222222222222222222"),
			Index:    0,
			Path:     "episode-01.mkv",
			Size:     2048,
		},
		{
			InfoHash: protocol.MustParseID("2222222222222222222222222222222222222222"),
			Index:    1,
			Path:     "episode-02.mkv",
			Size:     4096,
		},
	}
	fake := &fakeSearch{
		torrentFilesResult: search.TorrentFilesResult{
			Items: expected,
		},
	}
	resolver := torrentResolver{&Resolver{Search: fake}}

	got, err := resolver.Files(context.Background(), &model.Torrent{
		InfoHash:    expected[0].InfoHash,
		FilesStatus: model.FilesStatusMulti,
		FilesCount:  model.NewNullUint(uint(len(expected))),
	})
	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, 1, fake.torrentFilesCalls)
}

func TestTorrentFilesResolverLimitUsesSingleFileFallback(t *testing.T) {
	require.Equal(t, uint(1), torrentFilesResolverLimit(model.Torrent{
		FilesStatus: model.FilesStatusSingle,
	}))
}
