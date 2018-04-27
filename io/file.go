package io

import (
	"errors"
	"io/ioutil"
)

type FileReader struct{}

func (r *FileReader) Read(args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("A file is required")
	}

	filepath := args[0]
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
