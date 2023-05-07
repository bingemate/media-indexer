package controllers

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/gin-gonic/gin"
	"log"
)

type uploadResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func InitUploadController(engine *gin.RouterGroup, mediaUploader *features.MediaUploader) {
	engine.POST("/movie", func(c *gin.Context) {
		uploadMovie(c, mediaUploader)
	})
	engine.POST("/tv", func(c *gin.Context) {
		uploadTvShow(c, mediaUploader)
	})
}

// @Summary		Upload Movies
// @Description	Upload movies from the configured folder
// @Tags			Upload
// @Accept 		multipart/form-data
// @Param			upload[] formData file true "Files to upload"
// @Produce		json
// @Success		200	{object} uploadResponse
// @Failure		400	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router			/upload/movie [post]
func uploadMovie(c *gin.Context, mediaUploader *features.MediaUploader) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	files := form.File["upload[]"]
	if files == nil || len(files) == 0 {
		c.JSON(400, errorResponse{Error: "no files uploaded in upload[]"})
		return
	}
	log.Println("Uploading", len(files), "movies...")
	for _, file := range files {
		err := mediaUploader.UploadMovie(c, file)
		if err != nil {
			log.Println(err)
			c.JSON(500, errorResponse{Error: err.Error()})
			return
		}
	}
	c.JSON(200, uploadResponse{Message: "ok", Count: len(files)})
}

// @Summary		Upload TV Shows
// @Description	Upload TV Shows from the configured folder
// @Tags			Upload
// @Accept 		multipart/form-data
// @Param			upload[] formData file true "Files to upload"
// @Produce		json
// @Success		200	{object} uploadResponse
// @Failure		400	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router			/upload/tv [post]
func uploadTvShow(c *gin.Context, mediaUploader *features.MediaUploader) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	files := form.File["upload[]"]
	if files == nil || len(files) == 0 {
		c.JSON(400, errorResponse{Error: "no files uploaded in upload[]"})
		return
	}
	log.Println("Uploading", len(files), "tv shows...")
	for _, file := range files {
		err := mediaUploader.UploadTV(c, file)
		if err != nil {
			log.Println(err)
			c.JSON(500, errorResponse{Error: err.Error()})
			return
		}
	}
	c.JSON(200, uploadResponse{Message: "ok", Count: len(files)})
}
