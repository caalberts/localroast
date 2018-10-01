package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	argsFunc, assertArgs := setUpTestFunc()
	runEFunc, assertRunE := setUpTestFunc()

	cobraCmd := &cobra.Command{
		Args: cobra.PositionalArgs(argsFunc),
		RunE: runEFunc,
	}
	defaultCmd := &basicCommand{cobraCmd}
	rootCmd := newRootCmd(defaultCmd)

	cmd := rootCmd.getCommand()
	args := []string{"arg1", "arg2"}

	cmd.SetArgs(args)

	t.Run("validates using default command's validation", func(t *testing.T) {
		assert.NoError(t, cmd.ValidateArgs(args))
		assertArgs(t, true, args)
	})

	t.Run("executes default command", func(t *testing.T) {
		assert.NoError(t, cmd.Execute())
		assertRunE(t, true, args)
	})
}

type testFunc func(*cobra.Command, []string) error
type assertTestFunc func(*testing.T, bool, []string)

func setUpTestFunc() (testFunc, assertTestFunc) {
	var called bool
	var receivedArgs []string

	testFunc := func(cmd *cobra.Command, args []string) error {
		called = true
		receivedArgs = args
		return nil
	}

	assertTestFunc := func(t *testing.T, expectCalled bool, expectArgs []string) {
		assert.Equal(t, expectCalled, called)
		assert.Equal(t, expectArgs, receivedArgs)
	}

	return testFunc, assertTestFunc
}
