package pkg

import "os"

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
	return os.Rename(source, destination)
}
