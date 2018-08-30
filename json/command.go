package json

import (
	"io"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/spf13/afero"
	"log"
)

type validator interface {
	Validate([]string) error
}

type parser interface {
	Parse(io.Reader) ([]localroast.Schema, error)
}

// Command struct contains a file reader to read input file,
// a parser to parse input into schema,
// and a server constructor.
type Command struct {
	v  validator
	p  parser
	s  http.ServerFunc
	fs fileSystem
}

type fileSystem interface {
	Open(string) (afero.File, error)
}

// NewCommand creates a command with a JSON file reader and parser.
func NewCommand() Command {
	return Command{
		v:  Validator{},
		p:  Parser{},
		s:  http.NewServer,
		fs: afero.NewOsFs(),
	}
}

// Execute runs the command and start a server.
func (c Command) Execute(port string, args []string) error {
	server := c.s(port)

	if err := c.v.Validate(args); err != nil {
		return err
	}

	file, err := c.fs.Open(args[0])
	if err != nil {
		return err
	}

	schemas, err := c.p.Parse(file)
	if err != nil {
		return err
	}

	updateChan := server.Watch()
	updateChan <- schemas

	log.Println("brewing on port " + port)
	return server.ListenAndServe()
}
