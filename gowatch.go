package main

import (
	"errors"
	"fmt"
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
		Files chan string
	}
)

func main() {
	watcher, err := NewFileWatcher("./")
	if err != nil {
		log.Fatal(err)
	}
	go watcher.run()

	for {
		select {
		case <-watcher.Files:
			goTest()
		}
	}
}

// NewFileWatcher creates and returns a new FileWatcher
func NewFileWatcher(path string) (*FileWatcher, error) {
	folders := subfolders(path)
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

// run starts a goroutine to watch for changes on the FileWatcher
func (watcher *FileWatcher) run() {
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
func subfolders(path string) (paths []string) {
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

// goTest runs go test ./... and prints output
func goTest() {
	args := []string{"test", "./..."}
	cmd := exec.Command("go", args...)
	color.Cyan(strings.Join(cmd.Args, " "))

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	color.Yellow(string(out))

	var cmdState string
	if cmd.ProcessState.Success() {
		cmdState = color.GreenString("PASS")
	} else {
		cmdState = color.RedString("FAIL")
	}
	fmt.Println(cmdState, fmt.Sprintf("(%.2f seconds)", cmd.ProcessState.UserTime().Seconds()))
}
