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

type Movie struct {
	Type        int
	ID          int
	Name        string
	ReleaseDate string
}

func (m *Movie) Year() string {
	return strings.Split(m.ReleaseDate, "-")[0]
}

type AtomicMediaList struct {
	mediaList map[MediaFile]Movie
	lock      sync.Mutex
}

func (a *AtomicMediaList) LinkMediaFile(mediaFile MediaFile, media Movie) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.mediaList[mediaFile] = media
}

func (a *AtomicMediaList) Get(mediaFile MediaFile) (Movie, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	media, ok := a.mediaList[mediaFile]
	return media, ok
}

func (a *AtomicMediaList) GetAll() map[MediaFile]Movie {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.mediaList
}

func NewAtomicMediaList() *AtomicMediaList {
	return &AtomicMediaList{
		mediaList: make(map[MediaFile]Movie),
		lock:      sync.Mutex{},
	}
}

type MediaClient interface {
	SearchMovie(query string, year string) (Movie, error)
	SearchTVShow(query string, season, episode int) (Movie, error)
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
		Type:        MOVIE,
		ID:          media.ID,
		Name:        media.Title,
		ReleaseDate: media.ReleaseDate,
	}, nil
}

func (m *mediaClient) SearchTVShow(query string, season, episode int) (Movie, error) {
	//TODO implement me
	panic("implement me")
}
