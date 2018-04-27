package json

import (
	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/io"
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
}

func NewCommand() *Command {
	return &Command{
		r: &io.FileReader{},
		p: &Parser{},
	}
}

func (c *Command) Execute(args []string) ([]localroast.Schema, error) {
	bytes, err := c.r.Read(args)
	if err != nil {
		return nil, err
	}

	schema, err := c.p.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
