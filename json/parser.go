package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/caalberts/localroast"
)

type Parser struct{}

type stub struct {
	Method   *string         `json:"method"`
	Path     *string         `json:"path"`
	Status   *int            `json:"status"`
	Response json.RawMessage `json:"response"`
}

func (p Parser) Parse(in <-chan io.Reader, out chan<- []localroast.Schema) error {
	var stubs []stub
	if err := json.NewDecoder(<-in).Decode(&stubs); err != nil {
		return err
	}

	schemas := make([]localroast.Schema, len(stubs))
	for i, stub := range stubs {
		if f := missingFields(stub); len(f) > 0 {
			return fmt.Errorf("missing required fields: %s", strings.Join(f, ", "))
		}
		schemas[i] = localroast.Schema{
			Method:   *stub.Method,
			Path:     *stub.Path,
			Status:   *stub.Status,
			Response: []byte(stub.Response),
		}
	}

	out <- schemas
	return nil
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
