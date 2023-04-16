package controllers

import (
	"github.com/bingemate/media-indexer/initializers"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine, env initializers.Env) {
	InitScanController(engine.Group("/scan"), env)
}
