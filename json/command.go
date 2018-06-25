package json

import (
	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
)

type fileReader interface {
	Read([]string) ([]byte, error)
}

type parser interface {
	Parse([]byte) ([]localroast.Schema, error)
}

// Command struct contains a file reader to read input file,
// a parser to parse input into schema,
// and a server constructor.
type Command struct {
	r fileReader
	p parser
	s http.ServerFunc
}

// NewCommand creates a command with a JSON file reader and parser.
func NewCommand() Command {
	return Command{
		r: FileReader{},
		p: Parser{},
		s: http.NewServer,
	}
}

// Execute runs the command and start a server.
func (c Command) Execute(args []string) error {
	bytes, err := c.r.Read(args)
	if err != nil {
		return err
	}

	schema, err := c.p.Parse(bytes)
	if err != nil {
		return err
	}

	server := c.s("8080", schema)

	return server.ListenAndServe()
}
