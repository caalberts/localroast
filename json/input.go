package json

import (
	"errors"
	"strings"
)

type Validator struct{}

func (v Validator) Validate(args []string) error {
	if len(args) < 1 {
		return errors.New("A file is required")
	}

	if len(args) > 1 {
		return errors.New("Too many arguments")
	}

	file := args[0]
	if !strings.HasSuffix(file, ".json") {
		return errors.New("Input must be a JSON file")
	}

	return nil
}
