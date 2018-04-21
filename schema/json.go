package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/caalberts/localroast"
)

func FromJSON(filepath string) ([]localroast.Schema, error) {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []localroast.Schema{}, err
	}
	return BytesToSchema(f)
}

type stub struct {
	Method   *string         `json:"method"`
	Path     *string         `json:"path"`
	Status   *int            `json:"status"`
	Response json.RawMessage `json:"response"`
}

func BytesToSchema(bytes []byte) ([]localroast.Schema, error) {
	var stubs []stub
	err := json.Unmarshal(bytes, &stubs)
	if err != nil {
		return []localroast.Schema{}, err
	}

	schemas := make([]localroast.Schema, len(stubs))
	for i, stub := range stubs {
		if f := missingFields(stub); len(f) > 0 {
			return []localroast.Schema{}, fmt.Errorf("Missing required fields: %s", strings.Join(f, ", "))
		}
		schemas[i] = localroast.Schema{
			Method: *stub.Method,
			Path:   *stub.Path,
			Status: *stub.Status,
		}
	}

	return schemas, nil
}

func missingFields(s stub) []string {
	var missingFields []string
	if s.Method == nil {
		missingFields = append(missingFields, "method")
	}
	if s.Path == nil {
		missingFields = append(missingFields, "path")
	}
	if s.Status == nil {
		missingFields = append(missingFields, "status")
	}
	return missingFields
}
