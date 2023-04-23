package pkg

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func IsDirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

func MoveFile(source, destination string) error {
	log.Println("Opening source file:", source)
	srcFile, err := os.Open(source)
	if err != nil {
		log.Println("Error opening source file:", err)
		return err
	}
	defer srcFile.Close()

	destDir := filepath.Dir(destination)
	log.Println("Creating destination directory:", destDir)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Println("Error creating destination directory:", err)
		return err
	}

	log.Println("Creating destination file:", destination)
	destFile, err := os.Create(destination)
	if err != nil {
		log.Println("Error creating destination file:", err)
		return err
	}
	defer destFile.Close()

	log.Println("Copying file contents from source to destination")
	if _, err := io.Copy(destFile, srcFile); err != nil {
		log.Println("Error copying file contents:", err)
		return err
	}

	log.Println("Removing source file:", source)
	if err := os.Remove(source); err != nil {
		log.Println("Error removing source file:", err)
		return err
	}

	log.Println("File move completed successfully")
	return nil
}

func ClearFolderContent(source string) error {
	log.Println("Clearing folder content:", source)
	files, err := os.ReadDir(source)
	if err != nil {
		log.Println("Error clearing folder content:", err)
		return err
	}
	for _, f := range files {
		log.Println("Removing file:", f.Name())
		if err := os.RemoveAll(filepath.Join(source, f.Name())); err != nil {
			log.Println("Error removing file:", err)
			return err
		}
	}
	return nil
}
