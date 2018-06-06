package json

import (
	"errors"
	"testing"

	"github.com/caalberts/localroast"
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

func TestExecuteJSONCommand(t *testing.T) {
	r := new(mockReader)
	p := new(mockParser)

	args := []string{"filename"}
	mockResult := []byte("content")
	r.On("Read", args).Return(mockResult, nil)
	p.On("Parse", mockResult).Return([]localroast.Schema{}, nil)

	cmd := &Command{r, p}
	schema, err := cmd.Execute(args)

	assert.Nil(t, err)
	assert.Equal(t, []localroast.Schema{}, schema)

	r.AssertExpectations(t)
	p.AssertExpectations(t)
}

func TestReadError(t *testing.T) {
	r := new(mockReader)
	p := new(mockParser)

	args := []string{"fakefile"}
	errorMsg := "Failed to read file"
	r.On("Read", args).Return([]byte(""), errors.New(errorMsg))

	cmd := &Command{r, p}
	schema, err := cmd.Execute(args)

	assert.NotNil(t, err)
	assert.Equal(t, "Failed to read file", err.Error())
	assert.Nil(t, schema)

	r.AssertExpectations(t)
	p.AssertNotCalled(t, "Parse")
}

func TestParseError(t *testing.T) {
	mockResult := []byte("content")
	r := new(mockReader)
	p := new(mockParser)

	args := []string{"filename"}
	errorMsg := "Failed to parse schema"
	r.On("Read", args).Return(mockResult, nil)
	p.On("Parse", mockResult).Return([]localroast.Schema{}, errors.New(errorMsg))

	cmd := &Command{r, p}
	schema, err := cmd.Execute(args)

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())
	assert.Nil(t, schema)

	r.AssertExpectations(t)
	p.AssertExpectations(t)
}
