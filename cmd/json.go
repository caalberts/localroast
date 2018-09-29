package cmd

import (
	"github.com/caalberts/localroast/filesystem"
	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/json"
	"github.com/caalberts/localroast/types"
	"github.com/spf13/cobra"
	"io"

	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

func init() {
	rootCmd.AddCommand(jsonCmd)
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Use localroast with json file (default)",
	Long: `A tool to help developers stub external HTTP services quickly.
See https://github.com/caalberts/localroast/examples/stubs.json
for examples.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateJSON(args)
	},
	Example: "localroast json examples/stubs.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		fileHandler, err := filesystem.NewFileHandler()
		if err != nil {
			return err
		}

		command := command{
			fileHandler: fileHandler,
			parser:      json.NewParser(),
			serverFunc:  http.NewServer,
		}

		return command.execute(port, args)
	},
}

type fileHandler interface {
	Output() chan io.Reader
	Open(fileName string) error
	Watch() error
}

type parser interface {
	Output() chan []types.Schema
	Watch(chan io.Reader)
}

type command struct {
	fileHandler fileHandler
	parser      parser
	serverFunc  http.ServerFunc
}

func (c command) execute(port string, args []string) error {
	filepath := args[0]
	err := c.fileHandler.Open(filepath)
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

	log.Info("brewing on port " + port)
	return server.ListenAndServe()
}

func validateJSON(args []string) error {
	if len(args) < 1 {
		return errors.New("a file is required")
	}

	if len(args) > 1 {
		return errors.New("expected 1 argument")
	}

	file := args[0]
	if !strings.HasSuffix(file, ".json") {
		return errors.New("input must be a JSON file")
	}

	return nil
}
