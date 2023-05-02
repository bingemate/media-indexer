package controllers

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/gin-gonic/gin"
)

type movieScanResponse struct {
	Data []features.MovieScannerResult `json:"data"`
}

type tvScanResponse struct {
	Data []features.TVScannerResult `json:"data"`
}

func InitScanController(engine *gin.RouterGroup, movieScanner *features.MovieScanner, tvScanner *features.TVScanner) {
	engine.POST("/movie", func(c *gin.Context) {
		scanMovie(c, movieScanner)
	})
	engine.POST("/tv", func(c *gin.Context) {
		scanTvShow(c, tvScanner)
	})
}

// @Summary		Scan Movies
// @Description	Scan movies from the configured folder
// @Tags			Scan
// @Produce		json
// @Success		200	{object} movieScanResponse
// @Failure		500	{object} errorResponse
// @Router			/scan/movie [post]
func scanMovie(c *gin.Context, movieScanner *features.MovieScanner) {
	var result, err = movieScanner.ScanMovies()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, movieScanResponse{
		Data: *result,
	})
}

// @Summary		Scan TV Shows
// @Description	Scan TV Shows from the configured folder
// @Tags			Scan
// @Produce		json
// @Success		200	{object} tvScanResponse
// @Failure		500	{object} errorResponse
// @Router			/scan/tv [post]
func scanTvShow(c *gin.Context, tvScanner *features.TVScanner) {
	var result, err = tvScanner.ScanTV()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, tvScanResponse{
		Data: *result,
	})
}
