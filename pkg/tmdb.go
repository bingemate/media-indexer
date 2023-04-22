package pkg

import (
	"errors"
	"github.com/ryanbradynd05/go-tmdb"
	"strings"
	"sync"
)

const (
	TV_SHOW = iota
	TV_EPISODE
	MOVIE
)

type Media struct {
	Type        int
	ID          int
	Name        string
	ReleaseDate string
}

type Movie struct {
	Media
}

type TVShow struct {
	Media
}

type TVEpisode struct {
	Media
	TvShowID int
	Season   int
	Episode  int
}

func (m *Media) Year() string {
	return strings.Split(m.ReleaseDate, "-")[0]
}

type AtomicMediaList struct {
	mediaList map[MovieFile]Movie
	lock      sync.Mutex
}

func (a *AtomicMediaList) LinkMediaFile(mediaFile MovieFile, media Movie) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.mediaList[mediaFile] = media
}

func (a *AtomicMediaList) Get(mediaFile MovieFile) (Movie, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	media, ok := a.mediaList[mediaFile]
	return media, ok
}

func (a *AtomicMediaList) GetAll() map[MovieFile]Movie {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.mediaList
}

func NewAtomicMediaList() *AtomicMediaList {
	return &AtomicMediaList{
		mediaList: make(map[MovieFile]Movie),
		lock:      sync.Mutex{},
	}
}

type MediaClient interface {
	SearchMovie(query string, year string) (Movie, error)
	SearchTVShow(query string, season, episode int) (TVEpisode, error)
}

type mediaClient struct {
	tmdbClient *tmdb.TMDb
}

func NewMediaClient(apiKey string) MediaClient {
	config := tmdb.Config{
		APIKey:   apiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	return &mediaClient{
		tmdbClient: tmdb.Init(config),
	}
}

func (m *mediaClient) SearchMovie(query string, year string) (Movie, error) {
	var options = make(map[string]string)
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
	return Movie{
		Media: Media{
			Type:        MOVIE,
			ID:          media.ID,
			Name:        media.Title,
			ReleaseDate: media.ReleaseDate,
		},
	}, nil
}

func (m *mediaClient) SearchTVShow(query string, season, episode int) (TVEpisode, error) {
	//TODO implement me
	panic("implement me")
}
