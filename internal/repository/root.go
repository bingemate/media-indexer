package repository

import (
	"errors"
	"github.com/bingemate/media-go-pkg/repository"
	"github.com/bingemate/media-go-pkg/transcoder"
	"github.com/bingemate/media-indexer/pkg"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
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

	media := repository.Media{
		MediaType:   repository.MediaTypeMovie,
		TmdbID:      movie.ID,
		ReleaseDate: releaseDate,
	}
	inDb, err := r.findMedia(movie.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if inDb != nil {
		media.ID = inDb.ID
		media.CreatedAt = inDb.CreatedAt
	}
	db := r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}

	err = r.clearDuplicatedMovie(media.TmdbID, destinationPath, fileSource)
	if err != nil {
		return err
	}

	// Transcode movie here and retrieve file destination infos
	response, err := transcoder.ProcessFileTranscode(fileSource, media.ID, destinationPath, "15", "1280:720")
	if err != nil {
		return err
	}

	media.Categories = *r.extractCategories(&movie.Categories)
	media.Movies = []repository.Movie{
		{
			Name:      movie.Name,
			MediaFile: r.extractMediaFile(&mediaData, &response),
		},
	}

	db = r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) IndexTvEpisode(tvShow pkg.TVEpisode, fileSource, destinationPath string) error {
	log.Printf("Indexing tv show %s", tvShow.Name)
	releaseDate, err := time.Parse("2006-01-02", tvShow.TvReleaseDate)
	if err != nil {
		return err
	}
	episodeReleaseDate, err := time.Parse("2006-01-02", tvShow.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(fileSource)
	if err != nil {
		return err
	}
	tvShowEntity, err := r.handleTvShow(tvShow.Name, tvShow.TvShowID, releaseDate, &tvShow.Categories)
	if err != nil {
		return err
	}

	media := repository.Media{
		MediaType:   repository.MediaTypeEpisode,
		TmdbID:      tvShow.ID,
		ReleaseDate: episodeReleaseDate,
	}
	inDb, err := r.findMedia(tvShow.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if inDb != nil {
		media.ID = inDb.ID
		media.CreatedAt = inDb.CreatedAt
	}
	db := r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}

	err = r.clearDuplicatedEpisode(media.TmdbID, destinationPath, fileSource)
	if err != nil {
		return err
	}

	// Transcode episode here and retrieve file destination infos
	response, err := transcoder.ProcessFileTranscode(fileSource, media.ID, destinationPath, "15", "1280:720")
	if err != nil {
		return err
	}

	//media.Categories = *r.extractCategories(&tvShow.Categories)
	media.Episodes = []repository.Episode{
		{
			TvShow:    *tvShowEntity,
			Name:      tvShow.Name,
			NbEpisode: tvShow.Episode,
			NbSeason:  tvShow.Season,
			MediaFile: r.extractMediaFile(&mediaData, &response),
		},
	}

	db = r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) extractMediaFile(mediaData *pkg.MediaData, transcoderResponse *transcoder.TranscodeResponse) repository.MediaFile {
	return repository.MediaFile{
		Filename:  transcoderResponse.VideoIndex,
		Duration:  mediaData.Duration,
		Audio:     *r.extractAudio(&mediaData.Audios, transcoderResponse),
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
		alreadyInDB, err := r.getCategory(c.Name)
		if err != nil {
			categories[i] = repository.Category{
				Name: c.Name,
			}
		} else {
			categories[i] = *alreadyInDB
		}
	}
	return &categories
}

func (r *MediaRepository) clearDuplicatedMovie(tmdbID int, destination, fileDestination string) error {
	var movie repository.Movie
	db := r.db.Joins("Media").Joins("MediaFile").Where("tmdb_id = ?", tmdbID).First(&movie)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return db.Error
	}
	if movie.ID == "" {
		return nil
	}
	if movie.MediaFileID == "" {
		return r.db.Delete(&movie).Error
	}
	log.Printf("Removing duplicated movie %s", movie.Name)
	err := r.removeMediaFile(movie.MediaFileID)
	if err != nil {
		return err
	}
	if movie.MediaFile.Filename != fileDestination {
		log.Printf("Removing duplicated file %s", movie.MediaFile.Filename)
		err := os.RemoveAll(path.Join(destination, movie.MediaID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MediaRepository) clearDuplicatedEpisode(tmdbID int, destination, fileDestination string) error {
	var tvEpisode repository.Episode
	db := r.db.Joins("Media").Joins("MediaFile").Where("tmdb_id = ?", tmdbID).First(&tvEpisode)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return db.Error
	}
	if tvEpisode.ID == "" {
		return nil
	}
	if tvEpisode.MediaFileID == "" {
		return r.db.Delete(&tvEpisode).Error
	}
	log.Printf("Removing duplicated tv episode %s %dx%d", tvEpisode.Name, tvEpisode.NbSeason, tvEpisode.NbEpisode)
	err := r.removeMediaFile(tvEpisode.MediaFileID)
	if err != nil {
		return err
	}
	if tvEpisode.MediaFile.Filename != fileDestination {
		log.Printf("Removing duplicated file %s", tvEpisode.MediaFile.Filename)
		err := os.RemoveAll(path.Join(destination, tvEpisode.MediaID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MediaRepository) removeMediaFile(fileID string) error {
	return r.db.Delete(&repository.MediaFile{}, "id = ?", fileID).Error
}

func (r *MediaRepository) getCategory(name string) (*repository.Category, error) {
	var category repository.Category
	db := r.db.Where("name = ?", name).First(&category)
	if db.Error != nil {
		return nil, db.Error
	}
	return &category, nil
}

func (r *MediaRepository) handleTvShow(name string, tmdbID int, releaseDate time.Time, categories *[]pkg.Category) (*repository.TvShow, error) {
	var alreadyInDB repository.TvShow
	db := r.db.Joins("Media").Where("tmdb_id = ?", tmdbID).First(&alreadyInDB)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		log.Println(db.Error)
		return nil, db.Error
	}
	if alreadyInDB.ID != "" {
		return &alreadyInDB, nil
	}
	entity := &repository.TvShow{
		Media: repository.Media{
			MediaType:   repository.MediaTypeTvShow,
			TmdbID:      tmdbID,
			ReleaseDate: releaseDate,
			Categories:  *r.extractCategories(categories),
		},
		Name: name,
	}
	db = r.db.Save(entity)
	if db.Error != nil {
		return nil, db.Error
	}
	return entity, nil
}

func (r *MediaRepository) findMedia(tmdbID int) (*repository.Media, error) {
	var media repository.Media
	db := r.db.Where("tmdb_id = ?", tmdbID).First(&media)
	if db.Error != nil {
		return nil, db.Error
	}
	return &media, nil
}
