package video

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnrich(t *testing.T) {

	type parseTest struct {
		inputString    string
		contentType    model.ContentType
		expectedOutput model.TorrentContent
	}

	var parseTests = []parseTest{
		{
			inputString: "Mission.Impossible",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title: "Mission Impossible",
			},
		},
		{
			inputString: "Mission.Impossible.2023",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:       "Mission Impossible",
				ReleaseYear: 2023,
			},
		},
		{
			inputString: "Mission.Impossible.2023.1080p.BluRay.x264-SPARKS",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "Mission Impossible",
				ReleaseYear:     2023,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
				ReleaseGroup: model.NullString{
					String: "SPARKS",
					Valid:  true,
				},
			},
		},
		{
			inputString: "Die.Hard.(With.A.Vengeance!).And.A.Suffix.2023.1080p.BluRay.x264-SPARKS",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "Die Hard (With A Vengeance!) And A Suffix",
				ReleaseYear:     2023,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
				ReleaseGroup: model.NullString{
					String: "SPARKS",
					Valid:  true,
				},
			},
		},
		{
			inputString: "The.Movie.from.U.N.C.L.E.2015.1080p.BluRay.x264-SPARKS",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "The Movie from U.N.C.L.E.",
				ReleaseYear:     2015,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
				ReleaseGroup: model.NullString{
					String: "SPARKS",
					Valid:  true,
				},
			},
		},
		{
			inputString: "1776.1979.EXTENDED.HD.BluRay.X264-AMIABLE",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "1776",
				ReleaseYear:     1979,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
				ReleaseGroup: model.NullString{
					String: "AMIABLE",
					Valid:  true,
				},
			},
		},
		{
			inputString: "MY MOVIE (2016) [R][Action, Horror][720p.WEB-DL.AVC.8Bit.6ch.AC3].mkv",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "MY MOVIE",
				ReleaseYear:     2016,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBDL),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
			},
		},
		{
			inputString: "R.I.P.D.2013.720p.BluRay.x264-SPARKS",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "R.I.P.D.",
				ReleaseYear:     2013,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
				ReleaseGroup: model.NullString{
					String: "SPARKS",
					Valid:  true,
				},
			},
		},
		{
			inputString: "This Is A Movie (1999) [IMDB #] <Genre, Genre, Genre> {ACTORS} !DIRECTOR +MORE_SILLY_STUFF_NO_ONE_NEEDS ?",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:       "This Is A Movie",
				ReleaseYear: 1999,
			},
		},
		{
			inputString: "We Are the Movie!.2013.720p.H264.mkv",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "We Are the Movie!",
				ReleaseYear:     2013,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
			},
		},
		{
			inputString: "[ example.com ] We Are the Movie!.2013.720p.H264.mkv",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "We Are the Movie!",
				ReleaseYear:     2013,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
			},
		},
		{
			inputString: "Маша и Медведь в кино-12 месяцев.2022.WEBRip.1080p_от New-Team.mkv",
			contentType: model.ContentTypeMovie,
			expectedOutput: model.TorrentContent{
				Title:           "Маша и Медведь в кино-12 месяцев",
				ReleaseYear:     2022,
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBRip),
				//ReleaseGroup: "New-Team",
			},
		},
		{
			inputString: "The.Series.name.S04E08.1080p.WEB.h264-GRP[eztv.re].mkv",
			contentType: model.ContentTypeTvShow,
			expectedOutput: model.TorrentContent{
				Title:           "The Series name",
				Episodes:        make(model.Episodes).AddEpisode(4, 8),
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBRip),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				ReleaseGroup: model.NullString{
					String: "GRP",
					Valid:  true,
				},
			},
		},
		{
			inputString: "The.Series.name.S03-5.1080p.WEB.h264-GRP[eztv.re].mkv",
			contentType: model.ContentTypeTvShow,
			expectedOutput: model.TorrentContent{
				Title:           "The Series name",
				Episodes:        make(model.Episodes).AddSeason(3).AddSeason(4).AddSeason(5),
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBRip),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				ReleaseGroup: model.NullString{
					String: "GRP",
					Valid:  true,
				},
			},
		},
		{
			inputString: "The.Series.name.S04E03-5.1080p.WEB.h264-GRP[eztv.re].mkv",
			contentType: model.ContentTypeTvShow,
			expectedOutput: model.TorrentContent{
				Title:           "The Series name",
				Episodes:        make(model.Episodes).AddEpisode(4, 3).AddEpisode(4, 4).AddEpisode(4, 5),
				VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBRip),
				VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				ReleaseGroup: model.NullString{
					String: "GRP",
					Valid:  true,
				},
			},
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			torrent := model.Torrent{
				Name: test.inputString,
			}
			actualOutput, err := PreEnrich(model.TorrentContent{
				ContentType: model.NewNullContentType(test.contentType),
				Torrent:     torrent,
			})
			expectedOutput := test.expectedOutput
			expectedOutput.ContentType = model.NewNullContentType(test.contentType)
			expectedOutput.Torrent = torrent
			assert.NoError(t, err)
			assert.Equal(t, expectedOutput, actualOutput)
		})
	}
}
