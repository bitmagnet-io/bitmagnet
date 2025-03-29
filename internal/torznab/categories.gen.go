package torznab

var categoriesMap = map[int]Category{
	2000: {
		ID:   2000,
		Name: "Movies",
		Subcat: []Subcategory{
			{
				ID:   2030,
				Name: "Movies/SD",
			},
			{
				ID:   2040,
				Name: "Movies/HD",
			},
			{
				ID:   2045,
				Name: "Movies/UHD",
			},
			{
				ID:   2060,
				Name: "Movies/3D",
			},
		},
	},
	2030: {
		ID:     2030,
		Name:   "Movies/SD",
		Subcat: []Subcategory{},
	},
	2040: {
		ID:     2040,
		Name:   "Movies/HD",
		Subcat: []Subcategory{},
	},
	2045: {
		ID:     2045,
		Name:   "Movies/UHD",
		Subcat: []Subcategory{},
	},
	2060: {
		ID:     2060,
		Name:   "Movies/3D",
		Subcat: []Subcategory{},
	},
	3000: {
		ID:   3000,
		Name: "Audio",
		Subcat: []Subcategory{
			{
				ID:   3030,
				Name: "Audio/Audiobook",
			},
		},
	},
	3030: {
		ID:     3030,
		Name:   "Audio/Audiobook",
		Subcat: []Subcategory{},
	},
	4000: {
		ID:   4000,
		Name: "PC",
		Subcat: []Subcategory{
			{
				ID:   4050,
				Name: "PC/Games",
			},
		},
	},
	4050: {
		ID:     4050,
		Name:   "PC/Games",
		Subcat: []Subcategory{},
	},
	5000: {
		ID:   5000,
		Name: "TV",
		Subcat: []Subcategory{
			{
				ID:   5030,
				Name: "TV/SD",
			},
			{
				ID:   5040,
				Name: "TV/HD",
			},
			{
				ID:   5045,
				Name: "TV/UHD",
			},
		},
	},
	5030: {
		ID:     5030,
		Name:   "TV/SD",
		Subcat: []Subcategory{},
	},
	5040: {
		ID:     5040,
		Name:   "TV/HD",
		Subcat: []Subcategory{},
	},
	5045: {
		ID:     5045,
		Name:   "TV/UHD",
		Subcat: []Subcategory{},
	},
	6000: {
		ID:   6000,
		Name: "XXX",
		Subcat: []Subcategory{
			{
				ID:   6070,
				Name: "XXX/Other",
			},
		},
	},
	6070: {
		ID:     6070,
		Name:   "XXX/Other",
		Subcat: []Subcategory{},
	},
	7000: {
		ID:   7000,
		Name: "Books",
		Subcat: []Subcategory{
			{
				ID:   7020,
				Name: "Books/EBook",
			},
			{
				ID:   7030,
				Name: "Books/Comics",
			},
		},
	},
	7020: {
		ID:     7020,
		Name:   "Books/EBook",
		Subcat: []Subcategory{},
	},
	7030: {
		ID:     7030,
		Name:   "Books/Comics",
		Subcat: []Subcategory{},
	},
	8000: {
		ID:     8000,
		Name:   "Other",
		Subcat: []Subcategory{},
	},
}

var (
	CategoryMovies         = categoriesMap[2000]
	CategoryMoviesSD       = categoriesMap[2030]
	CategoryMoviesHD       = categoriesMap[2040]
	CategoryMoviesUHD      = categoriesMap[2045]
	CategoryMovies3D       = categoriesMap[2060]
	CategoryAudio          = categoriesMap[3000]
	CategoryAudioAudiobook = categoriesMap[3030]
	CategoryPC             = categoriesMap[4000]
	CategoryPCGames        = categoriesMap[4050]
	CategoryTV             = categoriesMap[5000]
	CategoryTVSD           = categoriesMap[5030]
	CategoryTVHD           = categoriesMap[5040]
	CategoryTVUHD          = categoriesMap[5045]
	CategoryXXX            = categoriesMap[6000]
	CategoryXXXOther       = categoriesMap[6070]
	CategoryBooks          = categoriesMap[7000]
	CategoryBooksEBook     = categoriesMap[7020]
	CategoryBooksComics    = categoriesMap[7030]
	CategoryOther          = categoriesMap[8000]
)

var TopLevelCategories = []Category{
	CategoryMovies,
	CategoryAudio,
	CategoryPC,
	CategoryTV,
	CategoryXXX,
	CategoryBooks,
	CategoryOther,
}
