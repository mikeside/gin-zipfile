package tool

import (
	"github.com/gookit/validate"
	"os"
)

func Mkdir(path string) (string, error) {

	if !validate.PathExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return path, err
		}
	}

	return path, nil
}
