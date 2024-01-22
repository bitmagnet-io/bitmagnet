package video

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {

	type output struct {
		contentType model.ContentType
		title       string
		releaseYear model.Year
		attrs       classifier.ContentAttributes
	}

	type parseTest struct {
		contentType    model.NullContentType
		inputString    string
		expectedOutput output
	}

	var parseTests = []parseTest{
		{
			inputString: "Mission.Impossible",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "Mission Impossible",
			},
		},
		{
			inputString: "Mission.Impossible.2023",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "Mission Impossible",
				releaseYear: 2023,
			},
		},
		{
			inputString: "Mission.Impossible.2023.1080p.BluRay.x264-SPARKS",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "Mission Impossible",
				releaseYear: 2023,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
					ReleaseGroup: model.NullString{
						String: "SPARKS",
						Valid:  true,
					},
				},
			},
		},
		{
			inputString: "Die.Hard.(With.A.Vengeance!).And.A.Suffix.2023.1080p.BluRay.x264-SPARKS",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "Die Hard (With A Vengeance!) And A Suffix",
				releaseYear: 2023,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
					ReleaseGroup: model.NullString{
						String: "SPARKS",
						Valid:  true,
					},
				},
			},
		},
		{
			inputString: "The.Movie.from.U.N.C.L.E.2015.1080p.BluRay.x264-SPARKS",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "The Movie from U.N.C.L.E.",
				releaseYear: 2015,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
					ReleaseGroup: model.NullString{
						String: "SPARKS",
						Valid:  true,
					},
				},
			},
		},
		{
			inputString: "1776.1979.EXTENDED.HD.BluRay.X264-AMIABLE",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "1776",
				releaseYear: 1979,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
					ReleaseGroup: model.NullString{
						String: "AMIABLE",
						Valid:  true,
					},
				},
			},
		},
		{
			inputString: "MY MOVIE (2016) [R][Action, Horror][720p.WEB-DL.AVC.8Bit.6ch.AC3].mkv",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "MY MOVIE",
				releaseYear: 2016,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBDL),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				},
			},
		},
		{
			inputString: "R.I.P.D.2013.720p.BluRay.x264-SPARKS",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "R.I.P.D.",
				releaseYear: 2013,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceBluRay),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecX264),
					ReleaseGroup: model.NullString{
						String: "SPARKS",
						Valid:  true,
					},
				},
			},
		},
		{
			inputString: "This Is A Movie (1999) [IMDB #] <Genre, Genre, Genre> {ACTORS} !DIRECTOR +MORE_SILLY_STUFF_NO_ONE_NEEDS ?",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "This Is A Movie",
				releaseYear: 1999,
			},
		},
		{
			inputString: "We Are the Movie!.2013.720p.H264.mkv",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "We Are the Movie!",
				releaseYear: 2013,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				},
			},
		},
		{
			inputString: "[ example.com ] We Are the Movie!.2013.720p.H264.mkv",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "We Are the Movie!",
				releaseYear: 2013,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV720p),
					VideoCodec:      model.NewNullVideoCodec(model.VideoCodecH264),
				},
			},
		},
		{
			inputString: "Маша и Медведь в кино-12 месяцев.2022.WEBRip.1080p_от New-Team.mkv",
			expectedOutput: output{
				contentType: model.ContentTypeMovie,
				title:       "Маша и Медведь в кино-12 месяцев",
				releaseYear: 2022,
				attrs: classifier.ContentAttributes{
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					VideoSource:     model.NewNullVideoSource(model.VideoSourceWEBRip),
					//ReleaseGroup: "New-Team",
				},
			},
		},
		{
			inputString: "The.Series.name.S04E08.1080p.WEB.h264-GRP[eztv.re].mkv",
			expectedOutput: output{
				contentType: model.ContentTypeTvShow,
				title:       "The Series name",
				attrs: classifier.ContentAttributes{
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
		},
		{
			inputString: "The.Series.name.S03-5.1080p.WEB.h264-GRP[eztv.re].mkv",
			expectedOutput: output{
				contentType: model.ContentTypeTvShow,
				title:       "The Series name",
				attrs: classifier.ContentAttributes{
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
		},
		{
			inputString: "The.Series.name.S04E03-5.1080p.WEB.h264-GRP[eztv.re].mkv",
			expectedOutput: output{
				contentType: model.ContentTypeTvShow,
				title:       "The Series name",
				attrs: classifier.ContentAttributes{
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
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			ct, title, year, attrs, err := ParseContent(
				test.contentType,
				test.inputString,
			)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutput.contentType, ct)
			assert.Equal(t, test.expectedOutput.title, title)
			assert.Equal(t, test.expectedOutput.releaseYear, year)
			assert.Equal(t, test.expectedOutput.attrs, attrs)
		})
	}
}
