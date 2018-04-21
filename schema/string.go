package schema

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/caalberts/localroast"
)

type String struct {
	Strings []string
}

func (s *String) CreateSchema() ([]localroast.Schema, error) {
	schemas := make([]localroast.Schema, len(s.Strings))
	for i, str := range s.Strings {
		schema, err := toSchema(str)
		if err != nil {
			return schemas, err
		}
		schemas[i] = schema
	}
	return schemas, nil
}

func toSchema(definition string) (localroast.Schema, error) {
	matches, err := validMatch(definition)
	if err != nil {
		return localroast.Schema{}, err
	}

	method := matches[1]
	path := matches[2]
	status, _ := strconv.Atoi(matches[3])
	schema := localroast.Schema{
		Method: method,
		Path:   path,
		Status: status,
	}
	return schema, nil
}

const regex = "^(GET|POST|PUT|PATCH|DELETE) ([\\w\\d/]+) (\\d{3})$"

var matcher = regexp.MustCompile(regex)

func validMatch(input string) ([]string, error) {
	matches := matcher.FindStringSubmatch(input)
	if len(matches) != 4 {
		return nil, errors.New("Invalid input: " + input)
	}

	return matches, nil
}
