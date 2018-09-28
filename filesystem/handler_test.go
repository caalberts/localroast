package filesystem

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const watchedFile = "watchthis.txt"

func TestMain(m *testing.M) {
	os.Create(watchedFile)

	testResult := m.Run()

	os.Remove(watchedFile)
	os.Exit(testResult)
}

func TestFileHandler_Add(t *testing.T) {
	eventChan := make(chan string)
	errorChan := make(chan error)

	fileHandler, err := NewFileHandler(eventChan, errorChan)
	assert.NoError(t, err)

	err = fileHandler.Add(watchedFile)
	assert.NoError(t, err)

	assert.Equal(t, watchedFile, fileHandler.file)
}

func TestFileHandler_Run(t *testing.T) {
	eventChan := make(chan string)
	errorChan := make(chan error)

	fileHandler, err := NewFileHandler(eventChan, errorChan)
	assert.NoError(t, err)

	err = fileHandler.Add(watchedFile)
	assert.NoError(t, err)

	fileHandler.Run()

	t.Run("sends event on WRITE", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Write,
		}
		fileHandler.watcher.Events <- event

		receivedFile := <-eventChan
		assert.Equal(t, watchedFile, receivedFile)
	})

	t.Run("sends event on CREATE", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Create,
		}
		fileHandler.watcher.Events <- event

		receivedFile := <-eventChan
		assert.Equal(t, watchedFile, receivedFile)
	})

	t.Run("sends error", func(t *testing.T) {
		fileHandler.watcher.Errors <- errors.New("something happened")

		receivedError := <-errorChan
		assert.Error(t, receivedError, "something happened")
	})
}
