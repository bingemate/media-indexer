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
	// Cors
	addCors(engine)

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

func addCors(engine *gin.Engine) gin.IRoutes {
	return engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func doc() {
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Media Indexer API"
	docs.SwaggerInfo.Description = "API for the Media Indexer"
}
