package repository

import (
	"errors"
	"fmt"
	"github.com/bingemate/media-go-pkg/repository"
	"github.com/bingemate/media-go-pkg/transcoder"
	"github.com/bingemate/media-indexer/pkg"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	if db == nil {
		log.Fatal("db is nil")
	}
	return &MediaRepository{db: db}
}

func (r *MediaRepository) IndexMovie(movie pkg.Movie, fileSource, destinationPath string) error {
	log.Printf("Indexing movie %s", movie.Name)
	releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(fileSource)
	if err != nil {
		return err
	}

	err = r.handleDuplicatedMovie(movie.ID, destinationPath)
	if err != nil {
		return err
	}

	// Transcode movie here and retrieve file destination infos
	response, err := transcoder.ProcessFileTranscode(fileSource, strconv.Itoa(movie.ID), destinationPath, "15", "1280:720")
	if err != nil {
		return err
	}

	alreadyInDB, err := r.findMovie(movie.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	movieEntity := repository.Movie{
		ID:          movie.ID,
		Name:        movie.Name,
		ReleaseDate: releaseDate,
		MediaFile:   r.extractMediaFile(&mediaData, &response),
	}

	if alreadyInDB != nil {
		movieEntity.CreatedAt = alreadyInDB.CreatedAt
	}

	db := r.db.Save(&movieEntity)
	if db.Error != nil {
		return db.Error
	}
	movieEntity.Categories = *r.extractCategories(&movie.Categories)
	db = r.db.Save(&movieEntity)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) IndexTvEpisode(tvEpisode pkg.TVEpisode, fileSource, destinationPath string) error {
	log.Printf("Indexing tv show %s", tvEpisode.Name)
	releaseDate, err := time.Parse("2006-01-02", tvEpisode.TvReleaseDate)
	if err != nil {
		return err
	}
	episodeReleaseDate, err := time.Parse("2006-01-02", tvEpisode.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(fileSource)
	if err != nil {
		return err
	}
	tvShowEntity, err := r.handleTvShow(tvEpisode.Name, tvEpisode.TvShowID, releaseDate, &tvEpisode.Categories)
	if err != nil {
		return err
	}

	err = r.handleDuplicatedEpisode(tvEpisode.ID, destinationPath)
	if err != nil {
		return err
	}

	// Transcode episode here and retrieve file destination infos
	response, err := transcoder.ProcessFileTranscode(fileSource, strconv.Itoa(tvEpisode.ID), destinationPath, "15", "1280:720")
	if err != nil {
		return err
	}

	alreadyInDB, err := r.findEpisode(tvEpisode.ID)
	if err != nil {
		return err
	}

	episodeEntity := repository.Episode{
		ID:          tvEpisode.ID,
		TvShow:      *tvShowEntity,
		Name:        fmt.Sprintf("%s %dx%02d", tvEpisode.Name, tvEpisode.Season, tvEpisode.Episode),
		NbEpisode:   tvEpisode.Episode,
		NbSeason:    tvEpisode.Season,
		ReleaseDate: episodeReleaseDate,
		MediaFile:   r.extractMediaFile(&mediaData, &response),
	}

	if alreadyInDB != nil {
		episodeEntity.CreatedAt = alreadyInDB.CreatedAt
	}

	db := r.db.Save(&episodeEntity)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) extractMediaFile(mediaData *pkg.MediaData, transcoderResponse *transcoder.TranscodeResponse) *repository.MediaFile {
	return &repository.MediaFile{
		Filename:  transcoderResponse.VideoIndex,
		Duration:  mediaData.Duration,
		Audios:    *r.extractAudio(&mediaData.Audios, transcoderResponse),
		Subtitles: *r.extractSubtitles(&mediaData.Subtitles, transcoderResponse),
	}
}

func (r *MediaRepository) extractSubtitles(pkgSubtitles *[]pkg.SubtitleData, transcoderResponse *transcoder.TranscodeResponse) *[]repository.Subtitle {
	var subtitles = make([]repository.Subtitle, len(*pkgSubtitles))
	for i, s := range *pkgSubtitles {
		subtitles[i] = repository.Subtitle{
			Filename: transcoderResponse.Subtitles[i].SubtitleIndex,
			Language: s.Language,
		}
	}
	return &subtitles
}

func (r *MediaRepository) extractAudio(audiosData *[]pkg.AudioData, transcoderResponse *transcoder.TranscodeResponse) *[]repository.Audio {
	var audio = make([]repository.Audio, len(*audiosData))
	for i, a := range *audiosData {
		audio[i] = repository.Audio{
			Filename: transcoderResponse.Audios[i].AudioIndex,
			Language: a.Language,
		}
	}
	return &audio
}

func (r *MediaRepository) extractCategories(pkgCategories *[]pkg.Category) *[]repository.Category {
	var categories = make([]repository.Category, len(*pkgCategories))
	for i, c := range *pkgCategories {
		InDB, err := r.getOrCreateCategory(c.Name)
		if err != nil {
			categories[i] = repository.Category{
				Name: c.Name,
			}
		} else {
			categories[i] = *InDB
		}
	}
	return &categories
}

func (r *MediaRepository) handleDuplicatedMovie(tmdbID int, destination string) error {
	var movie repository.Movie
	db := r.db.Joins("MediaFile").Where("movies.id = ?", tmdbID).First(&movie)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return db.Error
	}
	if movie.ID == 0 {
		return nil
	}

	if movie.MediaFileID != nil {
		log.Printf("Removing duplicated movie %s", movie.Name)
		err := r.removeMediaFile(*movie.MediaFileID)
		if err != nil {
			return err
		}
		log.Printf("Removing duplicated file %s", movie.MediaFile.Filename)
		return os.RemoveAll(path.Join(destination, strconv.Itoa(tmdbID)))
	}
	return nil
}

func (r *MediaRepository) handleDuplicatedEpisode(tmdbID int, destination string) error {
	var tvEpisode repository.Episode
	db := r.db.Joins("MediaFile").Where("episodes.id = ?", tmdbID).First(&tvEpisode)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return db.Error
	}
	if tvEpisode.ID == 0 {
		return nil
	}

	if tvEpisode.MediaFileID != nil {
		log.Printf("Removing duplicated tv episode %s %dx%d", tvEpisode.Name, tvEpisode.NbSeason, tvEpisode.NbEpisode)
		err := r.removeMediaFile(*tvEpisode.MediaFileID)
		if err != nil {
			return err
		}
		log.Printf("Removing duplicated file %s", tvEpisode.MediaFile.Filename)
		return os.RemoveAll(path.Join(destination, strconv.Itoa(tmdbID)))
	}
	return nil
}

func (r *MediaRepository) removeMediaFile(fileID string) error {
	return r.db.Delete(&repository.MediaFile{}, "id = ?", fileID).Error
}

func (r *MediaRepository) getOrCreateCategory(name string) (*repository.Category, error) {
	var category repository.Category
	db := r.db.Where("name = ?", name).First(&category)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, db.Error
	}
	if category.ID != "" {
		return &category, nil
	}
	category = repository.Category{
		Name: name,
	}
	db = r.db.Save(&category)
	if db.Error != nil {
		return nil, db.Error
	}
	return &category, nil
}

func (r *MediaRepository) handleTvShow(name string, tmdbID int, releaseDate time.Time, categories *[]pkg.Category) (*repository.TvShow, error) {
	var alreadyInDB repository.TvShow
	db := r.db.Where(`id = ?`, tmdbID).First(&alreadyInDB)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		log.Println(db.Error)
		return nil, db.Error
	}
	if alreadyInDB.ID != 0 {
		return &alreadyInDB, nil
	}
	entity := &repository.TvShow{
		ID:          tmdbID,
		ReleaseDate: releaseDate,
		Name:        name,
	}
	db = r.db.Save(entity)
	if db.Error != nil {
		return nil, db.Error
	}
	entity.Categories = *r.extractCategories(categories)
	db = r.db.Save(entity)
	if db.Error != nil {
		return nil, db.Error
	}
	return entity, nil
}

func (r *MediaRepository) findMovie(tmdbID int) (*repository.Movie, error) {
	var movie repository.Movie
	db := r.db.Where("id = ?", tmdbID).First(&movie)
	if db.Error != nil {
		return nil, db.Error
	}
	return &movie, nil
}

func (r *MediaRepository) findEpisode(tmdbID int) (*repository.Episode, error) {
	var episode repository.Episode
	db := r.db.Where("id = ?", tmdbID).First(&episode)
	if db.Error != nil {
		return nil, db.Error
	}
	return &episode, nil
}
