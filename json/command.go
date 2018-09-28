package json

import (
	"io"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/filesystem"
	"github.com/caalberts/localroast/http"
	"log"
)

type validator interface {
	Validate([]string) error
}

type fileHandler interface {
	Output() chan io.Reader
	Open(fileName string) error
	Watch() error
}

type parser interface {
	Output() chan []localroast.Schema
	Watch(chan io.Reader)
}

// Command struct contains a file reader to read input file,
// a parser to parse input into schema,
// and a server constructor.
type Command struct {
	validator   validator
	fileHandler fileHandler
	parser      parser
	serverFunc  http.ServerFunc
}

// NewCommand creates a command with a JSON file reader and parser.
func NewCommand() (*Command, error) {
	fileHandler, err := filesystem.NewFileHandler()
	if err != nil {
		return nil, err
	}

	cmd := Command{
		validator:   Validator{},
		fileHandler: fileHandler,
		parser:      NewParser(),
		serverFunc:  http.NewServer,
	}
	return &cmd, nil
}

// Execute runs the command and start a server.
func (c Command) Execute(port string, args []string) error {
	err := c.validator.Validate(args)
	if err != nil {
		return err
	}

	filepath := args[0]
	err = c.fileHandler.Open(filepath)
	if err != nil {
		return err
	}

	err = c.fileHandler.Watch()
	if err != nil {
		return err
	}

	server := c.serverFunc(port)

	c.parser.Watch(c.fileHandler.Output())
	server.Watch(c.parser.Output())

	log.Println("brewing on port " + port)
	return server.ListenAndServe()
}
