package io

import "errors"

type CLIReader struct{}

func (r *CLIReader) Read(args []string) ([]string, error) {
	if len(args) < 1 {
		return nil, errors.New("Please define an endpoint in the format '<METHOD> <PATH> <STATUS_CODE>'. e.g 'GET / 200'")
	}
	return args, nil
}
