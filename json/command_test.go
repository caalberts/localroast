package json

import (
	"errors"
	"io"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) Validate(strs []string) error {
	args := m.Called(strs)
	return args.Error(0)
}

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

func (m *mockParser) Output() chan []localroast.Schema {
	args := m.Called()
	return args.Get(0).(chan []localroast.Schema)
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

func (m *mockServer) Watch(schemas chan []localroast.Schema) {
	m.Called(schemas)
}

const (
	port     = "8080"
	testFile = "fixtures/test.json"
)

func resetMocks(v *mockValidator, f *mockFileHandler, p *mockParser, s *mockServer) {
	v.Mock = mock.Mock{}
	f.Mock = mock.Mock{}
	p.Mock = mock.Mock{}
	s.Mock = mock.Mock{}
}

func TestExecuteJSONCommand(t *testing.T) {
	mockValidator := new(mockValidator)
	mockFileHandler := new(mockFileHandler)
	mockParser := new(mockParser)
	mockServer := new(mockServer)

	var serverPort string
	sFunc := func(port string) http.Server {
		serverPort = port
		return mockServer
	}

	cmd := Command{mockValidator, mockFileHandler, mockParser, sFunc}

	t.Run("successful command", func(t *testing.T) {
		args := []string{testFile}

		readerChan := make(chan io.Reader)
		schemaChan := make(chan []localroast.Schema)

		mockValidator.On("Validate", args).Return(nil)
		mockFileHandler.On("Open", testFile).Return(nil)
		mockFileHandler.On("Watch").Return(nil)
		mockFileHandler.On("Output").Return(readerChan)
		mockParser.On("Watch", readerChan)
		mockParser.On("Output").Return(schemaChan)
		mockServer.On("ListenAndServe").Return(nil)
		mockServer.On("Watch", schemaChan)

		err := cmd.Execute(port, args)

		assert.NoError(t, err)
		assert.Equal(t, port, serverPort)

		mockValidator.AssertExpectations(t)
		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockValidator, mockFileHandler, mockParser, mockServer)
	})

	t.Run("invalid argument", func(t *testing.T) {
		args := []string{"fakefile"}
		errorMsg := "invalid argument"
		mockValidator.On("Validate", args).Return(errors.New(errorMsg))

		err := cmd.Execute(port, args)

		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())

		mockValidator.AssertExpectations(t)
		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockValidator, mockFileHandler, mockParser, mockServer)
	})

	t.Run("failed to open file", func(t *testing.T) {
		fileName := "missingfile"
		args := []string{fileName}
		errorMsg := "failed to open file"

		mockValidator.On("Validate", args).Return(nil)
		mockFileHandler.On("Open", fileName).Return(errors.New(errorMsg))

		err := cmd.Execute(port, args)

		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())

		mockValidator.AssertExpectations(t)
		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockValidator, mockFileHandler, mockParser, mockServer)
	})

	t.Run("failed to watch file", func(t *testing.T) {
		args := []string{testFile}
		errorMsg := "failed to watch file"

		mockValidator.On("Validate", args).Return(nil)
		mockFileHandler.On("Open", testFile).Return(nil)
		mockFileHandler.On("Watch").Return(errors.New(errorMsg))

		err := cmd.Execute(port, args)

		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())

		mockValidator.AssertExpectations(t)
		mockFileHandler.AssertExpectations(t)
		mockParser.AssertExpectations(t)
		mockServer.AssertExpectations(t)

		resetMocks(mockValidator, mockFileHandler, mockParser, mockServer)
	})
}
