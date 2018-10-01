package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version string
	port    string
)

func Execute(v string) {
	version = v
	newCommand().getCommand().Execute()
}

type commander interface {
	getCommand() *cobra.Command
}

type basicCommand struct {
	*cobra.Command
}

func (c *basicCommand) getCommand() *cobra.Command {
	return c.Command
}

func newCommand() commander {
	json := newJSONCmd()
	version := newVersionCmd()

	root := newRootCmd(json)
	root.getCommand().PersistentFlags().StringVarP(&port, "port", "p", "8080", "port number")

	addSubcommands(root, json, version)

	return root
}

func addSubcommands(parent commander, children ...commander) {
	parentCmd := parent.getCommand()
	for _, child := range children {
		parentCmd.AddCommand(child.getCommand())
	}
}

func newRootCmd(defaultCmder commander) commander {
	cmd := &cobra.Command{
		Use:   "localroast",
		Short: "Localroast quickly stubs a HTTP server",
		Long: `A tool to help developers stub external HTTP services quickly.
See https://github.com/caalberts/localroast/examples/stubs.json
for examples.`,
		Args:    defaultCmder.getCommand().Args,
		Example: "localroast examples/stubs.json",
		RunE:    defaultRunner(defaultCmder),
	}

	return &basicCommand{cmd}
}

func defaultRunner(defaultCmd commander) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		command := defaultCmd.getCommand()
		if err := command.Args(cmd, args); err != nil {
			return err
		}
		return command.RunE(cmd, args)
	}
}
