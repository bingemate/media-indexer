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
}

func scanMovie(c *gin.Context, movieScanner *features.MovieScanner) {
	var result, err = movieScanner.ScanMovieFolder()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"result": result,
	})
}

func scanTvShow(c *gin.Context, tvScanner *features.TVScanner) {
	var result, err = tvScanner.ScanTVFolder()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"result": result,
	})
}
