package cmd

import (
	"fmt"
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
	fmt.Println("Starting server on port", env.Port)
	err = engine.Run(":" + env.Port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(engine)
}
