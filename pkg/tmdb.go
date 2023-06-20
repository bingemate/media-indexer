package pkg

import (
	"errors"
	"github.com/bingemate/media-go-pkg/tmdb"
	"strings"
	"sync"
)

type Category struct {
	ID   int
	Name string
}

type Media struct {
	ID          int
	Name        string
	ReleaseDate string
	Categories  []Category
}

type Movie struct {
	Media
}

type TVEpisode struct {
	ID            int
	EpisodeName   string
	ReleaseDate   string
	Categories    []Category
	TvShowID      int
	TvShowName    string
	TvReleaseDate string
	Season        int
	Episode       int
}

func (m *Media) Year() string {
	return strings.Split(m.ReleaseDate, "-")[0]
}

func (m *TVEpisode) Year() string {
	return strings.Split(m.ReleaseDate, "-")[0]
}

type AtomicMovieList struct {
	mediaList map[MovieFile]Movie
	lock      sync.Mutex
}

type AtomicTVEpisodeList struct {
	mediaList map[TVShowFile]TVEpisode
	lock      sync.Mutex
}

func (a *AtomicMovieList) LinkMediaFile(mediaFile MovieFile, media Movie) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.mediaList[mediaFile] = media
}

func (a *AtomicMovieList) Get(mediaFile MovieFile) (Movie, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	media, ok := a.mediaList[mediaFile]
	return media, ok
}

func (a *AtomicMovieList) GetAll() map[MovieFile]Movie {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.mediaList
}

func NewAtomicMovieList() *AtomicMovieList {
	return &AtomicMovieList{
		mediaList: make(map[MovieFile]Movie),
		lock:      sync.Mutex{},
	}
}

func (a *AtomicTVEpisodeList) LinkMediaFile(mediaFile TVShowFile, media TVEpisode) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.mediaList[mediaFile] = media
}

func (a *AtomicTVEpisodeList) Get(mediaFile TVShowFile) (TVEpisode, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	media, ok := a.mediaList[mediaFile]
	return media, ok
}

func (a *AtomicTVEpisodeList) GetAll() map[TVShowFile]TVEpisode {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.mediaList
}

func NewAtomicTVEpisodeList() *AtomicTVEpisodeList {
	return &AtomicTVEpisodeList{
		mediaList: make(map[TVShowFile]TVEpisode),
		lock:      sync.Mutex{},
	}
}

type MediaClient interface {
	SearchMovie(query string, year string) (Movie, error)
	SearchTVShow(query string, season, episode int) (TVEpisode, error)
}

type mediaClient struct {
	client tmdb.MediaClient
}

func NewMediaClient(apiKey string) MediaClient {
	return &mediaClient{
		client: tmdb.NewMediaClient(apiKey),
	}
}

func NewRedisMediaClient(apiKey, redisHost, redisPassword string) MediaClient {
	return &mediaClient{
		client: tmdb.NewRedisMediaClient(apiKey, redisHost, redisPassword),
	}
}

func (m *mediaClient) SearchMovie(query string, year string) (Movie, error) {
	/*var options = make(map[string]string)
	options["language"] = "fr"
	if year != "" {
		options["year"] = year
	}
	results, err := m.tmdbClient.SearchMovie(query, options)
	if err != nil {
		return Movie{}, err
	}
	if results.TotalResults == 0 {
		return Movie{}, errors.New("no results found")
	}
	var media = results.Results[0]
	movieInfo, err := m.tmdbClient.GetMovieInfo(media.ID, options)
	if err != nil {
		return Movie{}, err
	}*/
	results, err := m.client.SearchMoviesYear(query, year, 1)
	if err != nil {
		return Movie{}, err
	}
	if results.TotalResult == 0 {
		return Movie{}, errors.New("no results found")
	}
	var media = results.Results[0]
	movieInfo, err := m.client.GetMovieShort(media.ID)

	var categories = make([]Category, 0)
	for _, genre := range movieInfo.Genres {
		categories = append(categories,
			Category{
				ID:   genre.ID,
				Name: genre.Name,
			})
	}

	return Movie{
		Media: Media{
			ID:          media.ID,
			Name:        media.Title,
			ReleaseDate: media.ReleaseDate,
			Categories:  categories,
		},
	}, nil
}

func (m *mediaClient) SearchTVShow(query string, season, episode int) (TVEpisode, error) {
	/*	var options = make(map[string]string)
		options["language"] = "fr"
		show, err := m.tmdbClient.SearchTv(query, options)
		if err != nil {
			return TVEpisode{}, err
		}
		if show.TotalResults == 0 {
			return TVEpisode{}, errors.New("no results found")
		}
		episodeInfo, err := m.tmdbClient.GetTvEpisodeInfo(show.Results[0].ID, season, episode, options)
		if err != nil {
			return TVEpisode{}, err
		}
		tvInfo, err := m.tmdbClient.GetTvInfo(show.Results[0].ID, options)
		if err != nil {
			return TVEpisode{}, err
		}*/
	tvShows, err := m.client.SearchTVShows(query, 1)
	if err != nil {
		return TVEpisode{}, err
	}
	if tvShows.TotalResult == 0 {
		return TVEpisode{}, errors.New("no results found")
	}
	tvShow := tvShows.Results[0]
	episodeInfo, err := m.client.GetTVEpisode(tvShow.ID, season, episode)
	if err != nil {
		return TVEpisode{}, err
	}

	var categories = make([]Category, 0)
	for _, genre := range tvShow.Genres {
		categories = append(categories,
			Category{
				ID:   genre.ID,
				Name: genre.Name,
			})
	}

	return TVEpisode{
		ID:            episodeInfo.ID,
		ReleaseDate:   episodeInfo.AirDate,
		EpisodeName:   episodeInfo.Name,
		Categories:    categories,
		TvShowName:    tvShow.Title,
		TvShowID:      tvShow.ID,
		TvReleaseDate: tvShow.ReleaseDate,
		Season:        season,
		Episode:       episode,
	}, nil
}
