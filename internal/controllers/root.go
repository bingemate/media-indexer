package controllers

import (
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(engine *gin.Engine, db *gorm.DB, env initializers.Env) {
	engine.MaxMultipartMemory = 32 << 20 // 32 MiB per file upload fragment
	var mediaClient = pkg.NewMediaClient(env.TMDBApiKey)
	var mediaRepository = repository.NewMediaRepository(db)
	var movieScanner = features.NewMovieScanner(env.MovieSourceFolder, env.MovieTargetFolder, mediaClient, mediaRepository)
	var tvScanner = features.NewTVScanner(env.TvSourceFolder, env.TvTargetFolder, mediaClient, mediaRepository)
	var mediaUploader = features.NewMediaUploader(env.TvSourceFolder, env.MovieSourceFolder)
	features.ScheduleScanner(env.ScanCron, movieScanner, tvScanner)
	InitScanController(engine.Group("/scan"), movieScanner, tvScanner)
	InitUploadController(engine.Group("/upload"), mediaUploader)
}
