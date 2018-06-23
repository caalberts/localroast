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

type Command struct {
	r fileReader
	p parser
	s http.ServerFunc
}

func NewCommand() Command {
	return Command{
		r: FileReader{},
		p: Parser{},
		s: http.NewServer,
	}
}

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
