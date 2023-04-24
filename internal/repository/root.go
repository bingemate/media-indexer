package repository

import (
	"errors"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/google/uuid"
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

func (r *MediaRepository) IndexMovie(movie *pkg.Movie, destination, fileDestination string) error {
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
				Name: movie.Name,
				MediaFile: MediaFile{
					Filename:  fileDestination,
					Size:      mediaData.Size,
					Duration:  mediaData.Duration,
					Codec:     VideoCodec(mediaData.Codec),
					Audio:     *r.extractAudio(&mediaData.Audios),
					Subtitles: *r.extractSubtitles(&mediaData.Subtitles),
				},
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

func (r *MediaRepository) IndexTvShow(tvShow *pkg.TVEpisode, destination, fileDestination string) error {
	log.Printf("Indexing tv show %s", tvShow.Name)
	releaseDate, err := time.Parse("2006-01-02", tvShow.TvReleaseDate)
	if err != nil {
		return err
	}
	epidodeReleaseDate, err := time.Parse("2006-01-02", tvShow.ReleaseDate)
	if err != nil {
		return err
	}
	mediaData, err := pkg.RetrieveMediaData(path.Join(destination, fileDestination))
	if err != nil {
		return err
	}

	media := Media{
		MediaType:   MediaTypeEpisode,
		TmdbID:      tvShow.ID,
		ReleaseDate: epidodeReleaseDate,
		Episodes: []Episode{
			{
				TvShow: TvShow{
					// TODO check if tv show already exists
					Media: Media{
						MediaType:   MediaTypeTvShow,
						TmdbID:      tvShow.TvShowID,
						ReleaseDate: releaseDate,
						Categories:  *r.extractCategories(&tvShow.Categories),
					},
					Name: tvShow.Name,
				},
				Name:      tvShow.Name,
				NbEpisode: tvShow.Episode,
				NbSeason:  tvShow.Season,
				MediaFile: MediaFile{
					Filename:  fileDestination,
					Size:      mediaData.Size,
					Duration:  mediaData.Duration,
					Codec:     VideoCodec(mediaData.Codec),
					Audio:     *r.extractAudio(&mediaData.Audios),
					Subtitles: *r.extractSubtitles(&mediaData.Subtitles),
				},
			},
		},
	}
	// TODO complete
	return nil
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
	if movie.ID == uuid.Nil {
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

func (r *MediaRepository) getCategory(name string) (*Category, error) {
	var category Category
	db := r.db.Where("name = ?", name).First(&category)
	if db.Error != nil {
		return nil, db.Error
	}
	return &category, nil
}
