package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/caalberts/localroast"
	log "github.com/sirupsen/logrus"
)

type Parser struct {
	output chan []localroast.Schema
}

func NewParser() *Parser {
	return &Parser{
		output: make(chan []localroast.Schema),
	}
}

func (p *Parser) Output() chan []localroast.Schema {
	return p.output
}

func (p *Parser) Watch(input chan io.Reader) {
	go func() {
		for input := range input {
			var stubs []stub
			err := json.NewDecoder(input).Decode(&stubs)
			if err != nil {
				log.Error(err)
			} else {
				schemas, err := createSchema(stubs)
				if err != nil {
					log.Error(err)
				} else {
					p.output <- schemas
				}
			}
		}
	}()
}

func createSchema(stubs []stub) ([]localroast.Schema, error) {
	schemas := make([]localroast.Schema, len(stubs))

	for i, stub := range stubs {
		if f := missingFields(stub); len(f) > 0 {
			err := fmt.Errorf("missing required fields: %s", strings.Join(f, ", "))
			return []localroast.Schema{}, err
		}
		schemas[i] = localroast.Schema{
			Method:   *stub.Method,
			Path:     *stub.Path,
			Status:   *stub.Status,
			Response: []byte(stub.Response),
		}
	}

	return schemas, nil
}

type stub struct {
	Method   *string         `json:"method"`
	Path     *string         `json:"path"`
	Status   *int            `json:"status"`
	Response json.RawMessage `json:"response"`
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
