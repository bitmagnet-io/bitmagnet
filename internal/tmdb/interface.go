package tmdb

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Client interface {
	ValidateApiKey(context.Context) error
	SearchMovie(context.Context, SearchMovieRequest) (SearchMovieResponse, error)
	MovieDetails(context.Context, MovieDetailsRequest) (MovieDetailsResponse, error)
	SearchTv(context.Context, SearchTvRequest) (SearchTvResponse, error)
	TvDetails(context.Context, TvDetailsRequest) (TvDetailsResponse, error)
	FindByID(context.Context, FindByIDRequest) (FindByIDResponse, error)
}

type SearchMovieRequest struct {
	Query              string
	IncludeAdult       bool
	Language           model.NullString
	PrimaryReleaseYear model.Year
	Year               model.Year
	Region             model.NullString
}

type SearchMovieResult struct {
	VoteCount        int64   `json:"vote_count"`
	ID               int64   `json:"id"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
	Title            string  `json:"title"`
	Popularity       float32 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIDs         []int64 `json:"genre_ids"`
	BackdropPath     string  `json:"backdrop_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}

type SearchMovieResponse struct {
	Page         int64               `json:"page"`
	TotalResults int64               `json:"total_results"`
	TotalPages   int64               `json:"total_pages"`
	Results      []SearchMovieResult `json:"results"`
}

type SearchTvRequest struct {
	Query            string
	IncludeAdult     bool
	Language         model.NullString
	FirstAirDateYear model.Year
	Year             model.Year
}

type SearchTvResult struct {
	OriginalName     string   `json:"original_name"`
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	VoteCount        int64    `json:"vote_count"`
	VoteAverage      float32  `json:"vote_average"`
	PosterPath       string   `json:"poster_path"`
	FirstAirDate     string   `json:"first_air_date"`
	Popularity       float32  `json:"popularity"`
	GenreIDs         []int64  `json:"genre_ids"`
	OriginalLanguage string   `json:"original_language"`
	BackdropPath     string   `json:"backdrop_path"`
	Overview         string   `json:"overview"`
	OriginCountry    []string `json:"origin_country"`
}

type SearchTvResponse struct {
	Page         int64            `json:"page"`
	TotalResults int64            `json:"total_results"`
	TotalPages   int64            `json:"total_pages"`
	Results      []SearchTvResult `json:"results"`
}

type MovieDetailsRequest struct {
	ID               int64
	AppendToResponse []string
	Language         model.NullString
}

type MovieDetailsResponse struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int64 `json:"budget"`
	Genres []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	ID                  int64   `json:"id"`
	IMDbID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float32 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		Name          string `json:"name"`
		ID            int64  `json:"id"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso3166_1 string `json:"iso_3166_1"`
		Name      string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int64  `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		Iso639_1 string `json:"iso_639_1"`
		Name     string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float32 `json:"vote_average"`
	VoteCount   int64   `json:"vote_count"`
}

type TvDetailsRequest struct {
	SeriesID         int64
	AppendToResponse []string
	Language         model.NullString
}

type TvDetailsResponse struct {
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		ID          int64  `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Gender      int    `json:"gender"`
		ProfilePath string `json:"profile_path"`
	} `json:"created_by"`
	EpisodeRunTime []int  `json:"episode_run_time"`
	FirstAirDate   string `json:"first_air_date"`
	Genres         []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string   `json:"homepage"`
	ID               int64    `json:"id"`
	InProduction     bool     `json:"in_production"`
	Languages        []string `json:"languages"`
	LastAirDate      string   `json:"last_air_date"`
	LastEpisodeToAir struct {
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ID             int64   `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		ProductionCode string  `json:"production_code"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int64   `json:"show_id"`
		StillPath      string  `json:"still_path"`
		VoteAverage    float32 `json:"vote_average"`
		VoteCount      int64   `json:"vote_count"`
	} `json:"last_episode_to_air"`
	Name             string `json:"name"`
	NextEpisodeToAir struct {
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ID             int64   `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		ProductionCode string  `json:"production_code"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int64   `json:"show_id"`
		StillPath      string  `json:"still_path"`
		VoteAverage    float32 `json:"vote_average"`
		VoteCount      int64   `json:"vote_count"`
	} `json:"next_episode_to_air"`
	Networks []struct {
		Name          string `json:"name"`
		ID            int64  `json:"id"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalName        string   `json:"original_name"`
	Overview            string   `json:"overview"`
	Popularity          float32  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		Name          string `json:"name"`
		ID            int64  `json:"id"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso3166_1 string `json:"iso_3166_1"`
		Name      string `json:"name"`
	} `json:"production_countries"`
	Seasons []struct {
		AirDate      string `json:"air_date"`
		EpisodeCount int    `json:"episode_count"`
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		PosterPath   string `json:"poster_path"`
		SeasonNumber int    `json:"season_number"`
	} `json:"seasons"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Type        string  `json:"type"`
	VoteAverage float32 `json:"vote_average"`
	VoteCount   int64   `json:"vote_count"`
	ExternalIDs struct {
		IMDbID      string `json:"imdb_id"`
		FreebaseMID string `json:"freebase_mid"`
		FreebaseID  string `json:"freebase_id"`
		TVDBID      int64  `json:"tvdb_id"`
		TVRageID    int64  `json:"tvrage_id"`
		FacebookID  string `json:"facebook_id"`
		InstagramID string `json:"instagram_id"`
		TwitterID   string `json:"twitter_id"`
		ID          int64  `json:"id,omitempty"`
	} `json:"external_ids,omitempty"`
}

type FindByIDRequest struct {
	ExternalSource string
	ExternalID     string
	Language       model.NullString
}

type FindByIDResponse struct {
	MovieResults []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIDs         []int64 `json:"genre_ids"`
		ID               int64   `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float32 `json:"vote_average"`
		VoteCount        int64   `json:"vote_count"`
		Popularity       float32 `json:"popularity"`
	} `json:"movie_results,omitempty"`
	TvResults []struct {
		OriginalName     string   `json:"original_name"`
		ID               int64    `json:"id"`
		Name             string   `json:"name"`
		VoteCount        int64    `json:"vote_count"`
		VoteAverage      float32  `json:"vote_average"`
		FirstAirDate     string   `json:"first_air_date"`
		PosterPath       string   `json:"poster_path"`
		GenreIDs         []int64  `json:"genre_ids"`
		OriginalLanguage string   `json:"original_language"`
		BackdropPath     string   `json:"backdrop_path"`
		Overview         string   `json:"overview"`
		OriginCountry    []string `json:"origin_country"`
		Popularity       float32  `json:"popularity"`
	} `json:"tv_results,omitempty"`
}
