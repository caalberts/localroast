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

func setup(fileContent string) afero.Fs {
	appFS := afero.NewMemMapFs()

	appFS.MkdirAll("fixtures", 0755)
	afero.WriteFile(appFS, testFile, []byte(fileContent), 0644)

	return appFS
}

func resetMocks(v *mockValidator, p *mockParser, s *mockServer) {
	v.Mock = mock.Mock{}
	p.Mock = mock.Mock{}
	s.Mock = mock.Mock{}
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
	fs := setup(validJSON)
	v := new(mockValidator)
	p := new(mockParser)
	s := new(mockServer)

	var parsedSchema []localroast.Schema
	var serverPort string
	sFunc := func(port string, schemas []localroast.Schema) http.Server {
		parsedSchema = schemas
		serverPort = port
		return s
	}

	cmd := Command{v, p, sFunc, fs}

	t.Run("successful command", func(t *testing.T) {
		args := []string{testFile}
		mockSchema := []localroast.Schema{
			{
				Method:   "GET",
				Path:     "/",
				Status:   200,
				Response: []byte("{}"),
			},
		}
		v.On("Validate", args).Return(nil)
		p.On("Parse", mock.Anything).Return(mockSchema, nil)
		s.On("ListenAndServe").Return(nil)

		err := cmd.Execute(port, args)

		assert.NoError(t, err)
		assert.Equal(t, mockSchema, parsedSchema)
		assert.Equal(t, port, serverPort)

		v.AssertExpectations(t)
		p.AssertExpectations(t)
		s.AssertExpectations(t)

		resetMocks(v, p, s)
	})

	t.Run("read error", func(t *testing.T) {
		args := []string{"fakefile"}
		errorMsg := "Failed to read file"
		v.On("Validate", args).Return(errors.New(errorMsg))

		err := cmd.Execute(port, args)

		assert.Error(t, err)
		assert.Equal(t, "Failed to read file", err.Error())

		v.AssertExpectations(t)
		p.AssertNotCalled(t, "Parse")
		s.AssertNotCalled(t, "ListenAndServe")

		resetMocks(v, p, s)
	})

	t.Run("parsing error", func(t *testing.T) {
		args := []string{testFile}
		errorMsg := "Failed to parse schema"
		v.On("Validate", args).Return(nil)
		p.On("Parse", mock.Anything).Return([]localroast.Schema{}, errors.New(errorMsg))

		err := cmd.Execute(port, args)

		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())

		v.AssertExpectations(t)
		p.AssertExpectations(t)
		s.AssertNotCalled(t, "ListenAndServe")

		resetMocks(v, p, s)
	})
}
