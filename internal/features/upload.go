package features

import (
	"fmt"
	"github.com/bingemate/media-indexer/pkg"
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
	locked := jobLock.TryLock()
	if !locked {
		log.Printf("Job '%s' already running, skipping this run", pkg.GetJobName())
		return fmt.Errorf("job '%s' already running, skipping this run", pkg.GetJobName())
	}
	defer jobLock.Unlock()
	pkg.ClearJobLogs("upload movie")
	pkg.AppendJobLog("Starting upload movie job")
	log.Println("Uploading movie", file.Filename)
	pkg.AppendJobLog(fmt.Sprintf("Uploading movie %s", file.Filename))
	return context.SaveUploadedFile(
		file,
		path.Join(m.movieSourceFolder, file.Filename),
	)
}

func (m *MediaUploader) UploadTV(context *gin.Context, file *multipart.FileHeader) error {
	locked := jobLock.TryLock()
	if !locked {
		log.Printf("Job '%s' already running, skipping this run", pkg.GetJobName())
		return fmt.Errorf("job '%s' already running, skipping this run", pkg.GetJobName())
	}
	defer jobLock.Unlock()
	pkg.ClearJobLogs("upload tv")
	pkg.AppendJobLog("Starting upload tv job")
	log.Println("Uploading TV", file.Filename)
	pkg.AppendJobLog(fmt.Sprintf("Uploading TV %s", file.Filename))
	return context.SaveUploadedFile(
		file,
		path.Join(m.tvSourceFolder, file.Filename),
	)
}
