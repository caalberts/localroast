package json

import (
	"errors"
	"io"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/spf13/afero"
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

type mockParser struct {
	mock.Mock
}

func (m *mockParser) Parse(r io.Reader) ([]localroast.Schema, error) {
	args := m.Called(r)
	return args.Get(0).([]localroast.Schema), args.Error(1)
}

type mockServer struct {
	mock.Mock
}

func (m *mockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

const (
	port     = "8080"
	testFile = "fixtures/test.json"
)

func setup(t *testing.T, fileContent string) afero.Fs {
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("fixtures", 0755)
	afero.WriteFile(appFS, testFile, []byte(fileContent), 0644)

	return appFS
}

func TestExecuteJSONCommand(t *testing.T) {
	validJSON := `
	[
		{
			"method": "GET",
			"path": "/",
			"status": 200,
			"response": {}
		}
	]`
	fs := setup(t, validJSON)
	r := new(mockValidator)
	p := new(mockParser)
	s := new(mockServer)

	var parsedSchema []localroast.Schema
	var serverPort string
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		parsedSchema = schemas
		serverPort = port
		return s
	}

	args := []string{testFile}
	mockSchema := []localroast.Schema{
		{
			Method:   "GET",
			Path:     "/",
			Status:   200,
			Response: []byte("{}"),
		},
	}
	r.On("Validate", args).Return(nil)
	p.On("Parse", mock.Anything).Return(mockSchema, nil)
	s.On("ListenAndServe").Return(nil)

	cmd := Command{r, p, sFunc, fs}
	err := cmd.Execute(port, args)

	assert.Nil(t, err)
	assert.Equal(t, mockSchema, parsedSchema)
	assert.Equal(t, port, serverPort)

	r.AssertExpectations(t)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
}

func TestReadError(t *testing.T) {
	fs := setup(t, "")
	r := new(mockValidator)
	p := new(mockParser)
	s := new(mockServer)
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		return s
	}

	args := []string{"fakefile"}
	errorMsg := "Failed to read file"
	r.On("Validate", args).Return(errors.New(errorMsg))

	cmd := Command{r, p, sFunc, fs}
	err := cmd.Execute(port, args)

	assert.NotNil(t, err)
	assert.Equal(t, "Failed to read file", err.Error())

	r.AssertExpectations(t)
	p.AssertNotCalled(t, "Parse")
	s.AssertNotCalled(t, "ListenAndServe")
}

func TestParseError(t *testing.T) {
	fs := setup(t, "")
	r := new(mockValidator)
	p := new(mockParser)
	s := new(mockServer)
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		return s
	}

	args := []string{testFile}
	errorMsg := "Failed to parse schema"
	r.On("Validate", args).Return(nil)
	p.On("Parse", mock.Anything).Return([]localroast.Schema{}, errors.New(errorMsg))

	cmd := Command{r, p, sFunc, fs}
	err := cmd.Execute(port, args)

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())

	r.AssertExpectations(t)
	p.AssertExpectations(t)
	s.AssertNotCalled(t, "ListenAndServe")
}
