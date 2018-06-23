package json

import (
	"errors"
	"io/ioutil"
	"strings"
)

type FileReader struct{}

func (r *FileReader) Read(args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("A file is required")
	}

	if len(args) > 1 {
		return nil, errors.New("Too many arguments")
	}

	file := args[0]
	if !strings.HasSuffix(file, ".json") {
		return nil, errors.New("Input must be a JSON file")
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
