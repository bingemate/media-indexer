package repository

import (
	"errors"
	"github.com/bingemate/media-indexer/pkg"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"time"
)

type (
	MediaType     string
	VideoCodec    string
	AudioCodec    string
	SubtitleCodec string
)

const (
	MediaTypeMovie   MediaType = "Movie"
	MediaTypeTvShow  MediaType = "TvShow"
	MediaTypeEpisode MediaType = "Episode"

	VideoCodecH264 VideoCodec = "H264"
	VideoCodecH265 VideoCodec = "H265"
	VideoCodecHEVC VideoCodec = "HEVC"

	AudioCodecAAC    AudioCodec = "AAC"
	AudioCodecAC3    AudioCodec = "AC3"
	AudioCodecEAC3   AudioCodec = "EAC3"
	AudioCodecMP3    AudioCodec = "MP3"
	AudioCodecDTS    AudioCodec = "DTS"
	AudioCodecVorbis AudioCodec = "Vorbis"

	SubtitleCodecSRT    SubtitleCodec = "SRT"
	SubtitleCodecSUBRIP SubtitleCodec = "SUBRIP"
	SubtitleCodecASS    SubtitleCodec = "ASS"
)

type Model struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MediaFile struct {
	Model
	Filename  string
	Size      float64
	Duration  float64
	Codec     VideoCodec
	Mimetype  string
	Audio     []Audio    `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
	Subtitles []Subtitle `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Media struct {
	Model
	MediaType   MediaType  `gorm:"index"`
	TmdbID      int        `gorm:"uniqueIndex"`
	ReleaseDate time.Time  `gorm:"type:date"`
	TvShows     []TvShow   `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Movies      []Movie    `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Episodes    []Episode  `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Categories  []Category `gorm:"many2many:category_media;constraint:OnDelete:CASCADE;"`
}

type TvShow struct {
	Model
	Name     string
	MediaID  string    `gorm:"type:uuid;not null"`
	Media    Media     `gorm:"reference:MediaID"`
	Episodes []Episode `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
}

type Episode struct {
	Model
	Name        string
	NbEpisode   int
	NbSeason    int
	MediaID     string    `gorm:"type:uuid;not null"`
	Media       Media     `gorm:"reference:MediaID"`
	TvShowID    string    `gorm:"type:uuid;not null"`
	TvShow      TvShow    `gorm:"reference:TvShowID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Movie struct {
	Model
	Name        string
	MediaID     string    `gorm:"type:uuid;not null"`
	Media       Media     `gorm:"reference:MediaID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Audio struct {
	Model
	Codec       AudioCodec
	Language    string
	Bitrate     float64
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Subtitle struct {
	Model
	Codec       SubtitleCodec
	Language    string
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Category struct {
	Model
	Name   string  `gorm:"uniqueIndex"`
	Medias []Media `gorm:"many2many:category_media"`
}

type CategoryMedia struct {
	MediaID    string   `gorm:"type:uuid;primaryKey"`
	Media      Media    `gorm:"reference:MediaID;constraint:OnDelete:CASCADE;"`
	CategoryID string   `gorm:"type:uuid;primaryKey"`
	Category   Category `gorm:"reference:CategoryID;constraint:OnDelete:CASCADE;"`
}

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
	media := Media{
		MediaType:   MediaTypeMovie,
		TmdbID:      movie.ID,
		ReleaseDate: releaseDate,
		Categories:  *r.extractCategories(&movie.Categories),
		Movies: []Movie{
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
	db := r.db.Create(&media)
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

	media := Media{
		MediaType:   MediaTypeEpisode,
		TmdbID:      tvShow.ID,
		ReleaseDate: episodeReleaseDate,
		Episodes: []Episode{
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
	db := r.db.Create(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) extractMediaFile(fileDestination string, mediaData *pkg.MediaData) MediaFile {
	return MediaFile{
		Filename:  fileDestination,
		Size:      mediaData.Size,
		Duration:  mediaData.Duration,
		Codec:     VideoCodec(mediaData.Codec),
		Audio:     *r.extractAudio(&mediaData.Audios),
		Subtitles: *r.extractSubtitles(&mediaData.Subtitles),
		Mimetype:  mediaData.Mimetype,
	}
}

func (r *MediaRepository) extractSubtitles(pkgSubtitles *[]pkg.SubtitleData) *[]Subtitle {
	var subtitles = make([]Subtitle, len(*pkgSubtitles))
	for i, s := range *pkgSubtitles {
		subtitles[i] = Subtitle{
			Codec:    SubtitleCodec(s.Codec),
			Language: s.Language,
		}
	}
	return &subtitles
}

func (r *MediaRepository) extractAudio(audiosData *[]pkg.AudioData) *[]Audio {
	var audio = make([]Audio, len(*audiosData))
	for i, a := range *audiosData {
		audio[i] = Audio{
			Codec:    AudioCodec(a.Codec),
			Language: a.Language,
			Bitrate:  a.Bitrate,
		}
	}
	return &audio
}

func (r *MediaRepository) extractCategories(pkgCategories *[]pkg.Category) *[]Category {
	var categories = make([]Category, len(*pkgCategories))
	for i, c := range *pkgCategories {
		alreadyInDB, err := r.getCategory(c.Name)
		if err != nil {
			categories[i] = Category{
				Name: c.Name,
			}
		} else {
			categories[i] = *alreadyInDB
		}
	}
	return &categories
}

func (r *MediaRepository) clearDuplicatedMovie(tmdbID int, destination, fileDestination string) error {
	var movie Movie
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
	var tvEpisode Episode
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
		err := tx.Delete(&Media{}, "id = ?", mediaID).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&MediaFile{}, "id = ?", fileID).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MediaRepository) removeEpisode(mediaID, fileID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&Media{}, "id = ?", mediaID).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&MediaFile{}, "id = ?", fileID).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MediaRepository) getCategory(name string) (*Category, error) {
	var category Category
	db := r.db.Where("name = ?", name).First(&category)
	if db.Error != nil {
		return nil, db.Error
	}
	return &category, nil
}

func (r *MediaRepository) handleTvShow(name string, tmdbID int, releaseDate time.Time, categories *[]pkg.Category) (*TvShow, error) {
	var alreadyInDB TvShow
	db := r.db.Joins("Media").Where("tmdb_id = ?", tmdbID).First(&alreadyInDB)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		log.Println(db.Error)
		return nil, db.Error
	}
	if alreadyInDB.ID != "" {
		return &alreadyInDB, nil
	}
	return &TvShow{
		Media: Media{
			MediaType:   MediaTypeTvShow,
			TmdbID:      tmdbID,
			ReleaseDate: releaseDate,
			Categories:  *r.extractCategories(categories),
		},
		Name: name,
	}, nil
}
