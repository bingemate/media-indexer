package repository

import (
	"errors"
	"github.com/bingemate/media-go-pkg/repository"
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

func (r *MediaRepository) IndexMovie(movie pkg.Movie, destination, fileDestination string) error {
	log.Printf("Indexing movie %s", movie.Name)
	releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(path.Join(destination, fileDestination))
	if err != nil {
		return err
	}
	media := repository.Media{
		MediaType:   repository.MediaTypeMovie,
		TmdbID:      movie.ID,
		ReleaseDate: releaseDate,
		Categories:  *r.extractCategories(&movie.Categories),
		Movies: []repository.Movie{
			{
				Name:      movie.Name,
				MediaFile: r.extractMediaFile(fileDestination, &mediaData),
			},
		},
	}
	err = r.clearDuplicatedMovie(movie.ID, destination, fileDestination)
	if err != nil {
		return err
	}
	inDb, err := r.findMedia(movie.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if inDb != nil {
		media.ID = inDb.ID
	}
	db := r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) IndexTvEpisode(tvShow pkg.TVEpisode, destination, fileDestination string) error {
	log.Printf("Indexing tv show %s", tvShow.Name)
	releaseDate, err := time.Parse("2006-01-02", tvShow.TvReleaseDate)
	if err != nil {
		return err
	}
	episodeReleaseDate, err := time.Parse("2006-01-02", tvShow.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(path.Join(destination, fileDestination))
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
		Episodes: []repository.Episode{
			{
				TvShow:    *tvShowEntity,
				Name:      tvShow.Name,
				NbEpisode: tvShow.Episode,
				NbSeason:  tvShow.Season,
				MediaFile: r.extractMediaFile(fileDestination, &mediaData),
			},
		},
	}
	err = r.clearDuplicatedEpisode(tvShow.ID, destination, fileDestination)
	if err != nil {
		return err
	}
	inDb, err := r.findMedia(tvShow.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if inDb != nil {
		media.ID = inDb.ID
	}
	db := r.db.Save(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) extractMediaFile(fileDestination string, mediaData *pkg.MediaData) repository.MediaFile {
	return repository.MediaFile{
		Filename:  fileDestination,
		Size:      mediaData.Size,
		Duration:  mediaData.Duration,
		Codec:     repository.VideoCodec(mediaData.Codec),
		Audio:     *r.extractAudio(&mediaData.Audios),
		Subtitles: *r.extractSubtitles(&mediaData.Subtitles),
		Mimetype:  mediaData.Mimetype,
	}
}

func (r *MediaRepository) extractSubtitles(pkgSubtitles *[]pkg.SubtitleData) *[]repository.Subtitle {
	var subtitles = make([]repository.Subtitle, len(*pkgSubtitles))
	for i, s := range *pkgSubtitles {
		subtitles[i] = repository.Subtitle{
			Codec:    repository.SubtitleCodec(s.Codec),
			Language: s.Language,
		}
	}
	return &subtitles
}

func (r *MediaRepository) extractAudio(audiosData *[]pkg.AudioData) *[]repository.Audio {
	var audio = make([]repository.Audio, len(*audiosData))
	for i, a := range *audiosData {
		audio[i] = repository.Audio{
			Codec:    repository.AudioCodec(a.Codec),
			Language: a.Language,
			Bitrate:  a.Bitrate,
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
	log.Printf("Removing duplicated movie %s", movie.Name)
	err := r.removeMovie(movie.MediaID, movie.MediaFileID)
	if err != nil {
		if movie.MediaFile.Filename != fileDestination {
			log.Printf("Removing duplicated file %s", movie.MediaFile.Filename)
			err := os.Remove(path.Join(destination, movie.MediaFile.Filename))
			if err != nil {
				return err
			}
		}
		return err
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
	log.Printf("Removing duplicated tv episode %s %dx%d", tvEpisode.Name, tvEpisode.NbSeason, tvEpisode.NbEpisode)
	err := r.removeEpisode(tvEpisode.MediaID, tvEpisode.MediaFileID)
	if err != nil {
		if tvEpisode.MediaFile.Filename != fileDestination {
			log.Printf("Removing duplicated file %s", tvEpisode.MediaFile.Filename)
			err := os.Remove(path.Join(destination, tvEpisode.MediaFile.Filename))
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

func (r *MediaRepository) removeMovie(mediaID, fileID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		//err := tx.Delete(&repository.Media{}, "id = ?", mediaID).Error
		//if err != nil {
		//	return err
		//}
		err := tx.Delete(&repository.MediaFile{}, "id = ?", fileID).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MediaRepository) removeEpisode(mediaID, fileID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		//err := tx.Delete(&repository.Media{}, "id = ?", mediaID).Error
		//if err != nil {
		//	return err
		//}
		err := tx.Delete(&repository.MediaFile{}, "id = ?", fileID).Error
		if err != nil {
			return err
		}
		return nil
	})
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
	return &repository.TvShow{
		Media: repository.Media{
			MediaType:   repository.MediaTypeTvShow,
			TmdbID:      tmdbID,
			ReleaseDate: releaseDate,
			Categories:  *r.extractCategories(categories),
		},
		Name: name,
	}, nil
}

func (r *MediaRepository) findMedia(tmdbID int) (*repository.Media, error) {
	var media repository.Media
	db := r.db.Where("tmdb_id = ?", tmdbID).First(&media)
	if db.Error != nil {
		return nil, db.Error
	}
	return &media, nil
}
