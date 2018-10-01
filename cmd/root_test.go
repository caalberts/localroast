package cmd

import (
	"io"
	"testing"

	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommandAndExecute(t *testing.T) {
	mockFileHandler := new(mockFileHandler)
	mockParser := new(mockParser)
	mockServer := new(mockServer)
	mockServerFunc := func(port string) http.Server {
		return mockServer
	}

	readerChan := make(chan io.Reader)
	schemaChan := make(chan []types.Schema)

	mockFileHandler.On("Open", "testfile.json").Return(nil)
	mockFileHandler.On("Watch").Return(nil)
	mockFileHandler.On("Output").Return(readerChan)
	mockParser.On("Watch", readerChan)
	mockParser.On("Output").Return(schemaChan)
	mockServer.On("ListenAndServe").Return(nil)
	mockServer.On("Watch", schemaChan)

	builder := &commandBuilder{
		fileHandler: mockFileHandler,
		jsonParser:  mockParser,
		serverFunc:  mockServerFunc,
	}
	cmd := builder.build()

	tests := []struct {
		name        string
		args        []string
		expectError string
	}{
		{"default command with json", []string{"testfile.json"}, ""},
		{"default command with non-json", []string{"testfile.txt"}, "must be a JSON file"},
		{"default command with too many args", []string{"blah.json", "testfile.json"}, "expected 1 argument"},
		{"json command with json", []string{"json", "testfile.json"}, ""},
		{"json command with non-json", []string{"json", "testfile.txt"}, "must be a JSON file"},
		{"json command with too many args", []string{"json", "blah.json", "testfile.json"}, "expected 1 argument"},
		{"unknown command", []string{"ruby"}, "must be a JSON file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.expectError)
			}
		})
	}
}

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
