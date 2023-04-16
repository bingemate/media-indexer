package cmd

import (
	"fmt"
	"github.com/bingemate/media-indexer/initializers"
)

func Serve(env initializers.Env) {
	fmt.Print(env)
}
