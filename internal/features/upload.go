package features

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"path"
)

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
	log.Println("Uploading movie", file.Filename)
	return context.SaveUploadedFile(
		file,
		path.Join(m.movieSourceFolder, file.Filename),
	)
}

func (m *MediaUploader) UploadTV(context *gin.Context, file *multipart.FileHeader) error {
	log.Println("Uploading TV", file.Filename)
	return context.SaveUploadedFile(
		file,
		path.Join(m.tvSourceFolder, file.Filename),
	)
}
