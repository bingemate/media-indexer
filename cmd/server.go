package cmd

import (
	"fmt"
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/controllers"
	"log"
)

func Serve(env initializers.Env) {
	var engine = initializers.InitGinEngine(env)
	controllers.InitRouter(engine, env)
	fmt.Println("Starting server on port", env.Port)
	err := engine.Run(":" + env.Port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(engine)
}
