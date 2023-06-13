package controllers

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/gin-gonic/gin"
)

type jobLogResponse pkg.JobLog

func InitJobController(engine *gin.RouterGroup) {
	engine.GET("/pop-logs", func(c *gin.Context) {
		popJobLogs(c)
	})
	engine.GET("/logs", func(c *gin.Context) {
		getJobLogs(c)
	})
	engine.GET("/job-name", func(c *gin.Context) {
		getJobName(c)
	})
	engine.GET("/is-running", func(c *gin.Context) {
		isRunning(c)
	})
}

// @Summary		Get Job Logs
// @Description	Get the logs of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{array} jobLogResponse
// @Router			/job/pop-logs [get]
func popJobLogs(c *gin.Context) {
	c.JSON(200, pkg.PopJobLogs())
}

// @Summary		Get Job Logs
// @Description	Get the logs of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{array} jobLogResponse
// @Router			/job/logs [get]
func getJobLogs(c *gin.Context) {
	c.JSON(200, pkg.GetJobLogs())
}

// @Summary		Get Job Name
// @Description	Get the name of the last / current job
// @Tags			Scan
// @Produce		json
// @Success		200	{string} string
// @Router			/job/job-name [get]
func getJobName(c *gin.Context) {
	c.JSON(200, pkg.GetJobName())
}

// @Summary		Is Running
// @Description	Check if a job is currently running
// @Tags			Scan
// @Produce		json
// @Success		200	{boolean} boolean
// @Router			/job/is-running [get]
func isRunning(c *gin.Context) {
	c.JSON(200, features.IsJobRunning())
}
