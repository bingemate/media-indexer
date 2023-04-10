package internal

import (
	"fmt"
	"github.com/bingemate/media-indexer/pkg/tree"
)

func Process(source, destination string) error {
	sourceTree, err := tree.BuildTree(source)
	if err != nil {
		return err
	}
	for _, mediaFile := range sourceTree {
		fmt.Println(mediaFile)
	}
	return nil
}
