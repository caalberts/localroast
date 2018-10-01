package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVersionCmd() commander {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print localroast version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("localroast %s", version)
		},
	}

	return &basicCommand{versionCmd}
}
