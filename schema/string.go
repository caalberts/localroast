package schema

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/caalberts/localroast"
)

const regex = "^(GET|POST) (/?[\\w/]+) (\\d{3})$"

var matcher = regexp.MustCompile(regex)

func FromString(definition string) (localroast.Schema, error) {
	matches, err := ValidMatch(definition)
	if err != nil {
		return localroast.Schema{}, err
	}

	method := matches[1]
	path := matches[2]
	code, _ := strconv.Atoi(matches[3])
	schema := localroast.Schema{
		Method:     method,
		Path:       path,
		StatusCode: code,
	}
	return schema, nil
}

func ValidMatch(input string) ([]string, error) {
	matches := matcher.FindStringSubmatch(input)
	if len(matches) != 4 {
		return nil, errors.New("Invalid input: " + input)
	}

	return matches, nil
}
