package controllers

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/gin-gonic/gin"
)

type jobLogResponse pkg.JobLog

func InitScanController(engine *gin.RouterGroup, movieScanner *features.MovieScanner, tvScanner *features.TVScanner) {
	engine.POST("/movie", func(c *gin.Context) {
		scanMovie(c, movieScanner)
	})
	engine.POST("/tv", func(c *gin.Context) {
		scanTvShow(c, tvScanner)
	})
	engine.GET("/pop-logs", func(c *gin.Context) {
		popJobLogs(c)
	})
	engine.GET("/logs", func(c *gin.Context) {
		getJobLogs(c)
	})
	engine.GET("/job-name", func(c *gin.Context) {
		getJobName(c)
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

// @Summary		Get Job Logs
// @Description	Get the logs of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{array} jobLogResponse
// @Router			/scan/pop-logs [get]
func popJobLogs(c *gin.Context) {
	c.JSON(200, pkg.PopJobLogs())
}

// @Summary		Get Job Logs
// @Description	Get the logs of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{array} jobLogResponse
// @Router			/scan/logs [get]
func getJobLogs(c *gin.Context) {
	c.JSON(200, pkg.GetJobLogs())
}

// @Summary		Get Job Name
// @Description	Get the name of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{string} string
// @Router			/scan/job-name [get]
func getJobName(c *gin.Context) {
	c.JSON(200, pkg.GetJobName())
}
