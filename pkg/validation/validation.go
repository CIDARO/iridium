package validation

import (
	"fmt"
	"os"
)

func ValidatePath(path string) (*string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", path)
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	return &path, nil
}
