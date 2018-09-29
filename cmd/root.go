package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version string
	port    string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "port number")
}

var rootCmd = &cobra.Command{
	Use:   "localroast",
	Short: "Localroast quickly stubs a HTTP server",
	Long: `A tool to help developers stub external HTTP services quickly.
See https://github.com/caalberts/localroast/examples/stubs.json
for examples.`,
	Args:    cobra.ExactArgs(1),
	Example: "localroast examples/stubs.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := jsonCmd.Args(cmd, args); err != nil {
			return err
		}
		return jsonCmd.RunE(cmd, args)
	},
}

func Execute(v string) {
	version = v
	rootCmd.Execute()
}
