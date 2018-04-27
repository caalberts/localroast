package strings

import (
	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/io"
)

type cliReader interface {
	Read([]string) ([]string, error)
}

type parser interface {
	Parse([]string) ([]localroast.Schema, error)
}

type Command struct {
	r cliReader
	p parser
}

func NewCommand() *Command {
	return &Command{
		r: &io.CLIReader{},
		p: &Parser{},
	}
}

func (c *Command) Execute(args []string) ([]localroast.Schema, error) {
	defs, err := c.r.Read(args)
	if err != nil {
		return nil, err
	}

	schema, err := c.p.Parse(defs)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
