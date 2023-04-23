package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MediaType string
type VideoCodec string
type AudioCodec string
type SubtitleCodec string

const (
	Movie   MediaType = "Movie"
	TvShow  MediaType = "TvShow"
	Episode MediaType = "Episode"

	H264 VideoCodec = "H264"
	H265 VideoCodec = "H265"

	AAC    AudioCodec = "AAC"
	AC3    AudioCodec = "AC3"
	MP3    AudioCodec = "MP3"
	DTS    AudioCodec = "DTS"
	Vorbis AudioCodec = "Vorbis"

	SRT    SubtitleCodec = "SRT"
	SUBRIP SubtitleCodec = "SUBRIP"
	ASS    SubtitleCodec = "ASS"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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

type Medias struct {
	Model
	MediaType   MediaType `gorm:"index"`
	TmdbID      int
	ReleaseDate time.Time    `gorm:"type:date"`
	TvShows     []TvShows    `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Movies      []Movies     `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Episodes    []Episodes   `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Categories  []Categories `gorm:"many2many:category_media"`
}

type TvShows struct {
	Model
	Name     string
	MediaID  string     `gorm:"type:uuid;not null"`
	Media    Medias     `gorm:"reference:MediaID"`
	Episodes []Episodes `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
}

type Episodes struct {
	Model
	Name        string
	NbEpisode   int
	NbSeason    int
	MediaID     string    `gorm:"type:uuid;not null"`
	Media       Medias    `gorm:"reference:MediaID"`
	TvShowID    string    `gorm:"type:uuid;not null"`
	TvShow      TvShows   `gorm:"reference:TvShowID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Movies struct {
	Model
	Name        string
	MediaID     string    `gorm:"type:uuid;not null"`
	Media       Medias    `gorm:"reference:MediaID"`
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

type Categories struct {
	Model
	Name   string
	Medias []Medias `gorm:"many2many:category_media"`
}

type CategoryMedia struct {
	MediaID    string     `gorm:"type:uuid;primaryKey"`
	Media      Medias     `gorm:"reference:MediaID"`
	CategoryID string     `gorm:"type:uuid;primaryKey"`
	Category   Categories `gorm:"reference:CategoryID"`
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}
