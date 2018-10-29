package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	fmt.Println(watcher)
}
