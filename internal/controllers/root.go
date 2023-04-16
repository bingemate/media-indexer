package controllers

import (
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine, env initializers.Env) {
	var movieScanner = features.NewMovieScanner(env.MovieSourceFolder, env.MovieTargetFolder, env.TMDBApiKey)
	InitScanController(engine.Group("/scan"), movieScanner)
}
