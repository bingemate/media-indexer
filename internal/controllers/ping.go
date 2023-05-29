package controllers

import "github.com/gin-gonic/gin"

// @Summary Ping
// @Description Ping
// @Tags Ping
// @Accept json
// @Produce json
// @Success 200 {object} string "pong"
// @Router /ping [get]
func InitPingController(gr *gin.RouterGroup) {
	gr.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
