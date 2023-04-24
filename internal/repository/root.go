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
