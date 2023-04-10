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

func (m *Media) Year() string {
	return strings.Split(m.ReleaseDate, "-")[0]
}

type AtomicMediaList struct {
	mediaList map[MediaFile]Media
	lock      sync.Mutex
}

func (a *AtomicMediaList) LinkMediaFile(mediaFile MediaFile, media Media) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.mediaList[mediaFile] = media
}

func (a *AtomicMediaList) Get(mediaFile MediaFile) (Media, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	media, ok := a.mediaList[mediaFile]
	return media, ok
}

func (a *AtomicMediaList) GetAll() map[MediaFile]Media {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.mediaList
}

func NewAtomicMediaList() *AtomicMediaList {
	return &AtomicMediaList{
		mediaList: make(map[MediaFile]Media),
		lock:      sync.Mutex{},
	}
}

type MediaClient interface {
	SearchMovie(query string, year string) (Media, error)
	SearchTVShow(query string, season, episode int) (Media, error)
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

func (m *mediaClient) SearchMovie(query string, year string) (Media, error) {
	var options = make(map[string]string)
	options["language"] = "fr"
	if year != "" {
		options["year"] = year
	}
	results, err := m.tmdbClient.SearchMovie(query, options)
	if err != nil {
		return Media{}, err
	}
	if results.TotalResults == 0 {
		return Media{}, errors.New("no results found")
	}
	var media = results.Results[0]
	return Media{
		Type:        MOVIE,
		ID:          media.ID,
		Name:        media.Title,
		ReleaseDate: media.ReleaseDate,
	}, nil
}

func (m *mediaClient) SearchTVShow(query string, season, episode int) (Media, error) {
	//TODO implement me
	panic("implement me")
}
