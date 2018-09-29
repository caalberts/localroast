package filesystem

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type FileHandler struct {
	output  chan io.Reader
	watcher *fsnotify.Watcher
	file    string
}

func NewFileHandler() (*FileHandler, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	output := make(chan io.Reader)
	fh := FileHandler{
		watcher: watcher,
		output:  output,
	}
	return &fh, nil
}

func (f *FileHandler) Output() chan io.Reader {
	return f.output
}

func (f *FileHandler) Open(fileName string) error {
	f.file = fileName
	file, err := os.Open(f.file)
	if err != nil {
		return err
	}

	go f.send(file)

	return nil
}

func (f *FileHandler) Watch() error {
	dir := filepath.Dir(f.file)
	err := f.watcher.Add(dir)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-f.watcher.Events:
				log.Debugf("received fsnotify event: %s", event)
				if f.writeOrCreate(event) && f.isWatched(event.Name) {
					log.Printf("file changed: %s", event.Name)
					file, err := os.Open(event.Name)
					if err != nil {
						log.Errorf("error opening file: %s", err)
					} else {
						f.send(file)
					}
				}
			case err := <-f.watcher.Errors:
				log.Errorf("error watching file: %s", err)
			}
		}
	}()

	return nil
}

func (f *FileHandler) send(file io.Reader) {
	f.output <- file
}

func (f *FileHandler) writeOrCreate(event fsnotify.Event) bool {
	return (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create)
}

func (f *FileHandler) isWatched(file string) bool {
	return file == f.file
}
