package watcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type (
	// FileWatcher is the struct used to watch files and folders for changes
	FileWatcher struct {
		*fsnotify.Watcher
		Files chan string
	}
)

// NewFileWatcher creates and returns a new FileWatcher
func NewFileWatcher(path string) (*FileWatcher, error) {
	folders, err := subfolders(path)
	if err != nil {
		return nil, err
	}

	if len(folders) == 0 {
		return nil, errors.New("no folders to watch")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &FileWatcher{Watcher: w}
	watcher.Files = make(chan string)

	for _, folder := range folders {
		watcher.addFolder(folder)
	}

	return watcher, nil
}

// Run watches for changes on the FileWatcher
func (watcher *FileWatcher) Run() {
	for {
		select {
		case event := <-watcher.Events:
			// create a file/directory
			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, err := os.Stat(event.Name)
				if err != nil {
					log.Println("error while processing file create", err)
				} else if fi.IsDir() {
					if !ShouldIgnoreFile(filepath.Base(event.Name)) {
						watcher.addFolder(event.Name)
					}
				} else {
					watcher.Files <- event.Name // created a file
				}
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				// modified a file, assuming that you don't modify folders
				watcher.Files <- event.Name
			}

		case err := <-watcher.Errors:
			log.Println("error while watching files", err)
		}
	}
}

// addFolder adds a folder to watch to the FileWatcher
func (watcher *FileWatcher) addFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}

	fmt.Println("Watching path: " + folder)
}

// subfolders returns a slice of subfolders including the current folder
func subfolders(path string) (paths []string, err error) {
	err = filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if ShouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}

		return nil
	})

	return paths, err
}

// ShouldIgnoreFile determines if a file should be ignored, file names that begin with "." or "_" are ignored by the go tool.
func ShouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || strings.HasPrefix(name, "vendor")
}
