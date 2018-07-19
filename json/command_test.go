package json

import (
	"errors"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Read(strs []string) ([]byte, error) {
	args := m.Called(strs)
	return args.Get(0).([]byte), args.Error(1)
}

type mockParser struct {
	mock.Mock
}

func (m *mockParser) Parse(bytes []byte) ([]localroast.Schema, error) {
	args := m.Called(bytes)
	return args.Get(0).([]localroast.Schema), args.Error(1)
}

type mockServer struct {
	mock.Mock
}

func (m *mockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

var port = "8080"

func TestExecuteJSONCommand(t *testing.T) {
	r := new(mockReader)
	p := new(mockParser)
	s := new(mockServer)

	var parsedSchema []localroast.Schema
	var serverPort string
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		parsedSchema = schemas
		serverPort = port
		return s
	}

	args := []string{"filename"}
	mockResult := []byte("content")
	mockSchema := []localroast.Schema{
		{
			Method:   "GET",
			Path:     "/",
			Status:   200,
			Response: []byte("{}"),
		},
	}
	r.On("Read", args).Return(mockResult, nil)
	p.On("Parse", mockResult).Return(mockSchema, nil)
	s.On("ListenAndServe").Return(nil)

	cmd := Command{r, p, sFunc}
	err := cmd.Execute(port, args)

	assert.Nil(t, err)
	assert.Equal(t, mockSchema, parsedSchema)
	assert.Equal(t, port, serverPort)

	r.AssertExpectations(t)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
}

func TestReadError(t *testing.T) {
	r := new(mockReader)
	p := new(mockParser)
	s := new(mockServer)
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		return s
	}

	args := []string{"fakefile"}
	errorMsg := "Failed to read file"
	r.On("Read", args).Return([]byte(""), errors.New(errorMsg))

	cmd := Command{r, p, sFunc}
	err := cmd.Execute(port, args)

	assert.NotNil(t, err)
	assert.Equal(t, "Failed to read file", err.Error())

	r.AssertExpectations(t)
	p.AssertNotCalled(t, "Parse")
	s.AssertNotCalled(t, "ListenAndServe")
}

func TestParseError(t *testing.T) {
	mockResult := []byte("content")
	r := new(mockReader)
	p := new(mockParser)
	s := new(mockServer)
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		return s
	}

	args := []string{"filename"}
	errorMsg := "Failed to parse schema"
	r.On("Read", args).Return(mockResult, nil)
	p.On("Parse", mockResult).Return([]localroast.Schema{}, errors.New(errorMsg))

	cmd := Command{r, p, sFunc}
	err := cmd.Execute(port, args)

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())

	r.AssertExpectations(t)
	p.AssertExpectations(t)
	s.AssertNotCalled(t, "ListenAndServe")
}
