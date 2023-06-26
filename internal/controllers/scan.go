package controllers

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/gin-gonic/gin"
)

func InitScanController(engine *gin.RouterGroup, movieScanner *features.MovieScanner, tvScanner *features.TVScanner) {
	engine.POST("/movie", func(c *gin.Context) {
		scanMovie(c, movieScanner)
	})
	engine.POST("/tv", func(c *gin.Context) {
		scanTvShow(c, tvScanner)
	})
	engine.POST("/all", func(c *gin.Context) {
		scanAll(c, movieScanner, tvScanner)
	})
}

// @Summary		Scan Movies
// @Description	Scan movies from the configured folder
// @Tags			Scan
// @Produce		json
// @Success		200	{string} string "Scan started"
// @Failure		500	{object} errorResponse
// @Router			/scan/movie [post]
func scanMovie(c *gin.Context, movieScanner *features.MovieScanner) {
	var err = movieScanner.ScanMovies()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, "Scan started")
}

// @Summary		Scan TV Shows
// @Description	Scan TV Shows from the configured folder
// @Tags			Scan
// @Produce		json
// @Success		200	{string} string "Scan started"
// @Failure		500	{object} errorResponse
// @Router			/scan/tv [post]
func scanTvShow(c *gin.Context, tvScanner *features.TVScanner) {
	var err = tvScanner.ScanTV()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, "Scan started")
}

// @Summary		Scan Movies and TV Shows
// @Description	Scan Movies and TV Shows from the configured folder
// @Tags			Scan
// @Produce		json
// @Success		200	{string} string "Scan started"
// @Failure		500	{object} errorResponse
// @Router			/scan/all [post]
func scanAll(c *gin.Context, movieScanner *features.MovieScanner, tvScanner *features.TVScanner) {
	var err = movieScanner.ScanMovies()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	err = tvScanner.ScanTV()
	if err != nil {
		c.JSON(500, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, "Scan started")
}
