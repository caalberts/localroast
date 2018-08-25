package filesystem

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFsnotifyWatcher_Run(t *testing.T) {
	watchedFile := "watchthis.txt"
	_, err := os.Create(watchedFile)
	assert.NoError(t, err)

	eventChan := make(chan string)
	errorChan := make(chan error)

	watcher, err := NewWatcher(eventChan, errorChan)
	assert.NoError(t, err)

	err = watcher.Add(watchedFile)
	assert.NoError(t, err)

	watcher.Run()

	t.Run("sends event on WRITE", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Write,
		}
		watcher.(*fsnotifyWatcher).watcher.Events <- event

		receivedFile := <-eventChan
		assert.Equal(t, watchedFile, receivedFile)
	})

	t.Run("sends event on CREATE", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Create,
		}
		watcher.(*fsnotifyWatcher).watcher.Events <- event

		receivedFile := <-eventChan
		assert.Equal(t, watchedFile, receivedFile)
	})

	t.Run("sends error", func(t *testing.T) {
		watcher.(*fsnotifyWatcher).watcher.Errors <- errors.New("something happened")

		receivedError := <-errorChan
		assert.Error(t, receivedError, "something happened")
	})

	err = os.Remove(watchedFile)
	assert.NoError(t, err)
}
