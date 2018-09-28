package filesystem

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
)

type FileHandler struct {
	watcher *fsnotify.Watcher
	file    string
	eventCh chan<- string
	errCh   chan<- error
}

func NewFileHandler(eventCh chan<- string, errCh chan<- error) (*FileHandler, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fh := FileHandler{
		watcher: watcher,
		eventCh: eventCh,
		errCh:   errCh,
	}
	return &fh, nil
}

func (f *FileHandler) Add(fileName string) error {
	if err := f.watcher.Add(fileName); err != nil {
		return err
	}

	f.file = fileName
	return nil
}

func (f *FileHandler) Run() {
	go func() {
		for {
			select {
			case event := <-f.watcher.Events:
				if f.writeOrCreate(event) && f.isWatched(event.Name) {
					f.eventCh <- event.Name
				}
			case err := <-f.watcher.Errors:
				f.errCh <- err
			}
		}
	}()
}

func (f *FileHandler) writeOrCreate(event fsnotify.Event) bool {
	return (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create)
}

func (f *FileHandler) isWatched(filePath string) bool {
	_, file := filepath.Split(filePath)
	return file == f.file
}
