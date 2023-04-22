package pkg

import (
	"errors"
	"github.com/ryanbradynd05/go-tmdb"
	"strings"
	"sync"
)

type Media struct {
	ID          int
	Name        string
	ReleaseDate string
}

type Movie struct {
	Media
}

type TVEpisode struct {
	Media
	TvShowID      int
	TvReleaseDate string
	Season        int
	Episode       int
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
			ID:          media.ID,
			Name:        media.Title,
			ReleaseDate: media.ReleaseDate,
		},
	}, nil
}

func (m *mediaClient) SearchTVShow(query string, season, episode int) (TVEpisode, error) {
	var options = make(map[string]string)
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
	return TVEpisode{
		Media: Media{
			ID:          episodeInfo.ID,
			ReleaseDate: episodeInfo.AirDate,
			Name:        show.Results[0].Name,
		},
		TvShowID:      show.Results[0].ID,
		TvReleaseDate: show.Results[0].FirstAirDate,
		Season:        season,
		Episode:       episode,
	}, nil
}
