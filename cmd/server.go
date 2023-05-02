package cmd

import (
	"fmt"
	"github.com/bingemate/media-indexer/docs"
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

func Serve(env initializers.Env) {
	var engine = gin.Default()
	db, err := initializers.ConnectToDB(env)
	if err != nil {
		log.Fatal(err)
	}
	controllers.InitRouter(engine, db, env)
	doc()
	fmt.Println("Starting server on port", env.Port)
	err = engine.Run(":" + env.Port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(engine)
}

func doc() {
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Media Indexer API"
	docs.SwaggerInfo.Description = "API for the Media Indexer"
}
