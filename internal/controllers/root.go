package controllers

import (
	objectstorage "github.com/bingemate/media-go-pkg/object-storage"
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type errorResponse struct {
	Error string `json:"error"`
}

func InitRouter(engine *gin.Engine, db *gorm.DB, env initializers.Env) {
	var mediaIndexerGroup = engine.Group("/media-indexer")
	engine.MaxMultipartMemory = 32 << 20 // 32 MiB per file upload fragment
	var mediaClient = pkg.NewRedisMediaClient(env.TMDBApiKey, env.RedisHost, env.RedisPassword)
	var mediaRepository = repository.NewMediaRepository(db, env.IntroFilePath)
	objectStorage, err := objectstorage.NewObjectStorage(env.S3AccessKeyId, env.S3SecretAccessKey, env.S3Endpoint, "fr-par", env.S3BucketName)
	if err != nil {
		panic(err)
	}
	var movieScanner = features.NewMovieScanner(env.MovieSourceFolder, env.MovieTargetFolder, mediaClient, mediaRepository, objectStorage)
	var tvScanner = features.NewTVScanner(env.TvSourceFolder, env.TvTargetFolder, mediaClient, mediaRepository, objectStorage)
	var mediaUploader = features.NewMediaUploader(env.TvSourceFolder, env.MovieSourceFolder)
	features.ScheduleScanner(env.ScanCron, movieScanner, tvScanner)
	InitScanController(mediaIndexerGroup.Group("/scan"), movieScanner, tvScanner)
	InitUploadController(mediaIndexerGroup.Group("/upload"), mediaUploader)
	InitJobController(mediaIndexerGroup.Group("/job"))
	InitPingController(mediaIndexerGroup.Group("/ping"))
}
