package filesystem

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	watchedFile = "watchthis.txt"
	fileContent = "hello world"
)

func TestMain(m *testing.M) {
	file, _ := os.Create(watchedFile)
	file.Write([]byte(fileContent))
	file.Close()

	testResult := m.Run()

	os.Remove(watchedFile)
	os.Exit(testResult)
}

func TestFileHandler_Open(t *testing.T) {
	fileHandler, err := NewFileHandler()
	assert.NoError(t, err)

	t.Run("with valid file", func(t *testing.T) {
		err = fileHandler.Open(watchedFile)
		assert.NoError(t, err)

		assert.Equal(t, watchedFile, fileHandler.file)

		result, err := ioutil.ReadAll(<-fileHandler.Output())
		assert.NoError(t, err)
		assert.Equal(t, []byte(fileContent), result)
	})

	t.Run("with non-existent file", func(t *testing.T) {
		assert.Error(t, fileHandler.Open("missing_file"))
	})
}

func TestFileHandler_Watch(t *testing.T) {
	fileHandler, err := NewFileHandler()
	assert.NoError(t, err)

	fileHandler.file = watchedFile
	err = fileHandler.Watch()
	assert.NoError(t, err)

	t.Run("sends file on WRITE as io.Reader", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Write,
		}
		fileHandler.watcher.Events <- event

		result, err := ioutil.ReadAll(<-fileHandler.Output())
		assert.NoError(t, err)
		assert.Equal(t, []byte(fileContent), result)
	})

	t.Run("sends file on CREATE as io.Reader", func(t *testing.T) {
		event := fsnotify.Event{
			Name: watchedFile,
			Op:   fsnotify.Create,
		}
		fileHandler.watcher.Events <- event

		result, err := ioutil.ReadAll(<-fileHandler.Output())
		assert.NoError(t, err)
		assert.Equal(t, []byte(fileContent), result)
	})

	t.Run("does not block when filewatch errors", func(t *testing.T) {
		fileHandler.watcher.Errors <- errors.New("something happened")

		assert.Empty(t, fileHandler.Output())
	})
}
