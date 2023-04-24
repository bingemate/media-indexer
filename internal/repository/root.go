package repository

import (
	"errors"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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

	AudioCodecAAC    AudioCodec = "AAC"
	AudioCodecAC3    AudioCodec = "AC3"
	AudioCodecMP3    AudioCodec = "MP3"
	AudioCodecDTS    AudioCodec = "DTS"
	AudioCodecVorbis AudioCodec = "Vorbis"

	SubtitleCodecSRT    SubtitleCodec = "SRT"
	SubtitleCodecSUBRIP SubtitleCodec = "SUBRIP"
	SubtitleCodecASS    SubtitleCodec = "ASS"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
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
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Movie struct {
	Model
	Name        string
	MediaID     string    `gorm:"type:uuid;not null"`
	Media       Media     `gorm:"reference:MediaID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
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

func (r *MediaRepository) IndexMovie(movie *pkg.Movie, fileDestination string) error {
	releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
	if err != nil {
		return err
	}
	mediadata, err := pkg.RetrieveMediaData(fileDestination)
	if err != nil {
		return err
	}
	media := Media{
		MediaType:   MediaTypeMovie,
		TmdbID:      movie.ID,
		ReleaseDate: releaseDate,
		Categories: func() []Category {
			var categories []Category
			for _, c := range movie.Categories {
				categories = append(categories, Category{
					Name: c.Name,
				})
			}
			return categories
		}(),
		Movies: []Movie{
			{
				Name: movie.Name,
				MediaFile: MediaFile{
					Filename: fileDestination,
					Size:     mediadata.Size,
					Duration: mediadata.Duration,
					Codec:    VideoCodec(mediadata.Codec),
					Audio: func() []Audio {
						var audio = make([]Audio, 0)
						for _, a := range mediadata.Audios {
							audio = append(audio, Audio{
								Codec:    AudioCodec(a.Codec),
								Language: a.Language,
								Bitrate:  a.Bitrate,
							})
						}
						return audio
					}(),
					Subtitles: func() []Subtitle {
						var subtitles = make([]Subtitle, 0)
						for _, s := range mediadata.Subtitles {
							subtitles = append(subtitles, Subtitle{
								Codec:    SubtitleCodec(s.Codec),
								Language: s.Language,
							})
						}
						return subtitles
					}(),
				},
			},
		},
	}
	err = r.clearDuplicatedMovie(movie.ID)
	if err != nil {
		return err
	}
	db := r.db.Create(&media)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *MediaRepository) clearDuplicatedMovie(tmdbID int) error {
	var movie Movie
	db := r.db.Joins("Media").Joins("MediaFile").Where("tmdb_id = ?", tmdbID).First(&movie)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return db.Error
	}
	if movie.ID == uuid.Nil {
		return nil
	}
	//TODO Si nom du fichier différent, on supprime l'ancien

	err := r.removeMovie(movie.MediaID, movie.MediaFileID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MediaRepository) removeMovie(mediaID, fileID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Delete(&Media{}, mediaID)
		if tx.Error != nil {
			return tx.Error
		}
		tx = tx.Delete(&MediaFile{}, fileID)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	})
}
