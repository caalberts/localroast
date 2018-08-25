package filesystem

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
)

type Watcher interface {
	Add(fileName string) error
	Run()
}

type fsnotifyWatcher struct {
	watcher *fsnotify.Watcher
	file    string
	eventCh chan<- string
	errCh   chan<- error
}

func NewWatcher(eventCh chan<- string, errCh chan<- error) (Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := fsnotifyWatcher{
		watcher: fsWatcher,
		eventCh: eventCh,
		errCh:   errCh,
	}
	return &watcher, nil
}

func (w *fsnotifyWatcher) Add(fileName string) error {
	if err := w.watcher.Add(fileName); err != nil {
		return err
	}

	w.file = fileName
	return nil
}

func (w *fsnotifyWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				if w.writeOrCreate(event) && w.isWatched(event.Name) {
					w.eventCh <- event.Name
				}
			case err := <-w.watcher.Errors:
				w.errCh <- err
			}
		}
	}()
}

func (w *fsnotifyWatcher) writeOrCreate(event fsnotify.Event) bool {
	return (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create)
}

func (w *fsnotifyWatcher) isWatched(filePath string) bool {
	_, file := filepath.Split(filePath)
	return file == w.file
}
