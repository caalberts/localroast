package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFileHandler struct {
	mock.Mock
}

func (m *mockFileHandler) Output() chan io.Reader {
	args := m.Called()
	return args.Get(0).(chan io.Reader)
}

func (m *mockFileHandler) Open(fileName string) error {
	args := m.Called(fileName)
	return args.Error(0)
}

func (m *mockFileHandler) Watch() error {
	args := m.Called()
	return args.Error(0)
}

type mockParser struct {
	mock.Mock
}

func (m *mockParser) Output() chan []types.Schema {
	args := m.Called()
	return args.Get(0).(chan []types.Schema)
}

func (m *mockParser) Watch(reader chan io.Reader) {
	m.Called(reader)
}

type mockServer struct {
	mock.Mock
}

func (m *mockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServer) Watch(schemas chan []types.Schema) {
	m.Called(schemas)
}

const (
	testPort = "8080"
	testFile = "fixtures/test.json"
)

func TestExecuteJSONCommand(t *testing.T) {
	mockFileHandler := new(mockFileHandler)
	mockParser := new(mockParser)
	mockServer := new(mockServer)

	var serverPort string
	sFunc := func(port string) http.Server {
		serverPort = port
		return mockServer
	}

	jsonCmd := newJSONCmd(mockFileHandler, mockParser, sFunc)
	cmd := jsonCmd.getCommand()
	cmd.Flags().String("port", testPort, "")

	t.Run("successful command", func(t *testing.T) {
		args := []string{testFile}

		readerChan := make(chan io.Reader)
		schemaChan := make(chan []types.Schema)

		mockFileHandler.On("Open", testFile).Return(nil)
		mockFileHandler.On("Watch").Return(nil)
		mockFileHandler.On("Output").Return(readerChan)
		mockParser.On("Watch", readerChan)
		mockParser.On("Output").Return(schemaChan)
		mockServer.On("ListenAndServe").Return(nil)
		mockServer.On("Watch", schemaChan)

		cmd.SetArgs(args)
		err := cmd.Execute()

		assert.NoError(t, err)
		assert.Equal(t, testPort, serverPort)

		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockFileHandler, mockParser, mockServer)
	})

	t.Run("failed to open file", func(t *testing.T) {
		fileName := "missingfile.json"
		args := []string{fileName}
		errorMsg := "failed to open file"

		mockFileHandler.On("Open", fileName).Return(errors.New(errorMsg))

		cmd.SetArgs(args)
		err := cmd.Execute()

		assert.Error(t, err)
		assert.Contains(t, errorMsg, err.Error())

		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockFileHandler, mockParser, mockServer)
	})

	t.Run("failed to watch file", func(t *testing.T) {
		args := []string{testFile}
		errorMsg := "failed to watch file"

		mockFileHandler.On("Open", testFile).Return(nil)
		mockFileHandler.On("Watch").Return(errors.New(errorMsg))

		cmd.SetArgs(args)
		err := cmd.Execute()

		assert.Error(t, err)
		assert.Contains(t, errorMsg, err.Error())

		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockFileHandler, mockParser, mockServer)
	})
}

func resetMocks(f *mockFileHandler, p *mockParser, s *mockServer) {
	f.Mock = mock.Mock{}
	p.Mock = mock.Mock{}
	s.Mock = mock.Mock{}
}

func TestValidateJSONArgs(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		err := validateJSONArgs(&cobra.Command{}, []string{"../examples/stubs.json"})
		assert.Nil(t, err)
	})

	t.Run("non json file", func(t *testing.T) {
		err := validateJSONArgs(&cobra.Command{}, []string{"stubs.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "input must be a JSON file", err.Error())
	})

	t.Run("without argument", func(t *testing.T) {
		err := validateJSONArgs(&cobra.Command{}, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "a file is required", err.Error())
	})

	t.Run("too many arguments", func(t *testing.T) {
		err := validateJSONArgs(&cobra.Command{}, []string{"abc.json", "def.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "expected 1 argument", err.Error())
	})

	t.Run("with incorrect file", func(t *testing.T) {
		err := validateJSONArgs(&cobra.Command{}, []string{"unknownfile"})
		assert.NotNil(t, err)
	})
}
