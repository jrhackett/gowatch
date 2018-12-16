package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type (
	// FileWatcher is the struct used to watch files and folders for changes
	FileWatcher struct {
		*fsnotify.Watcher
		Files   chan string
		Folders chan string
	}
)

func main() {
	watcher, err := NewFileWatcher("./")
	if err != nil {
		panic(err)
	}
	go watcher.Run()

	for {
		select {
		case <-watcher.Files:
			args := []string{"test", "./..."}
			cmd := exec.Command("go", args...)
			color.Yellow(strings.Join(cmd.Args, " "))

			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
			}
			color.Cyan(string(out))

			if cmd.ProcessState.Success() {
				color.Green("PASS")
			} else {
				color.Red("FAIL")
			}
			color.Yellow(" (%.2f seconds)\n", cmd.ProcessState.UserTime().Seconds())
		case folder := <-watcher.Folders:
			color.Yellow("Watching path: " + folder)
		}
	}
}

// NewFileWatcher creates and returns a new FileWatcher
func NewFileWatcher(path string) (*FileWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("no folders to watch")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher := &FileWatcher{Watcher: w}

	watcher.Files = make(chan string, 10)
	watcher.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		watcher.AddFolder(folder)
	}
	return watcher, nil
}

// Run starts a goroutine to watch for changes on the FileWatcher
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
					if !shouldIgnoreFile(filepath.Base(event.Name)) {
						watcher.AddFolder(event.Name)
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

// AddFolder adds a folder to watch to the FileWatcher
func (watcher *FileWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	watcher.Folders <- folder
}

// Subfolders returns a slice of subfolders including the current folder
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored, file names that begin with "." or "_" are ignored by the go tool.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || strings.HasPrefix(name, "vendor")
}
