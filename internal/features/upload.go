package features

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"path"
	"sync"
)

var uploadLock = sync.Mutex{}

type MediaUploader struct {
	tvSourceFolder    string
	movieSourceFolder string
}

func NewMediaUploader(tvSourceFolder, movieSourceFolder string) *MediaUploader {
	return &MediaUploader{
		tvSourceFolder:    tvSourceFolder,
		movieSourceFolder: movieSourceFolder,
	}
}

func (m *MediaUploader) UploadMovie(context *gin.Context, file *multipart.FileHeader) error {
	locked := uploadLock.TryLock()
	if !locked {
		log.Println("Upload or scan is already in progress")
		return errors.New("upload or scan is already in progress")
	}
	defer uploadLock.Unlock()
	log.Println("Uploading movie", file.Filename)
	return context.SaveUploadedFile(
		file,
		path.Join(m.movieSourceFolder, file.Filename),
	)
}

func (m *MediaUploader) UploadTV(context *gin.Context, file *multipart.FileHeader) error {
	locked := uploadLock.TryLock()
	if !locked {
		log.Println("Upload or scan is already in progress")
		return errors.New("upload or scan is already in progress")
	}
	defer uploadLock.Unlock()
	log.Println("Uploading TV", file.Filename)
	return context.SaveUploadedFile(
		file,
		path.Join(m.tvSourceFolder, file.Filename),
	)
}
