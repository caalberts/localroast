package cmd

import (
	"errors"
	"github.com/caalberts/localroast/http"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

type jsonCommand struct {
	*cobra.Command

	fileHandler fileHandler
	parser      parser
	serverFunc  http.ServerFunc
}

func (c *jsonCommand) getCommand() *cobra.Command {
	return c.Command
}

func newJSONCmd(fileHandler fileHandler, parser parser, serverFunc http.ServerFunc) commander {
	jsonCmd := &jsonCommand{
		fileHandler: fileHandler,
		parser:      parser,
		serverFunc:  serverFunc,
	}

	command := &cobra.Command{
		Use:   "json",
		Short: "Use localroast with json file (default)",
		Long: `A tool to help developers stub external HTTP services quickly.
	See https://github.com/caalberts/localroast/examples/stubs.json
	for examples.`,
		Args:    validateJSONArgs,
		Example: "localroast json examples/stubs.json",
		RunE:    jsonCmd.execute,
	}

	jsonCmd.Command = command

	return jsonCmd
}

func (c *jsonCommand) execute(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetString("port")
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

	log.Info("brewing on port " + port)
	return server.ListenAndServe()
}

func validateJSONArgs(cmd *cobra.Command, args []string) error {
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
