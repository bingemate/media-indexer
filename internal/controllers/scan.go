package controllers

import (
	"github.com/gin-gonic/gin"
)

func InitScanController(engine *gin.RouterGroup) {
	engine.POST("/movie", scanMovie)
	engine.POST("/tvshow", scanTvShow)
}

func scanMovie(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func scanTvShow(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
