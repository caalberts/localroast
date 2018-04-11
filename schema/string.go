package schema

import (
	"regexp"
	"strconv"

	"github.com/caalberts/localghost"
)

func FromString(definition string) (localghost.Schema, error) {
	regex := "^(GET) (/) (\\d{3})$"
	matches := regexp.MustCompile(regex).FindStringSubmatch(definition)

	method := matches[1]
	path := matches[2]
	code, _ := strconv.Atoi(matches[3])
	schema := localghost.Schema{
		Method:     method,
		Path:       path,
		StatusCode: code,
	}
	return schema, nil
}
