package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/jrhackett/gowatch/internal/watcher"
)

func main() {
	watcher, err := watcher.NewFileWatcher("./")
	if err != nil {
		log.Fatal(err)
	}

	go watcher.Run()

	for {
		<-watcher.Files
		goTest()
	}
}

// goTest runs go test ./... and prints output
func goTest() {
	cmd := exec.Command("go", []string{"test", "./..."}...)
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
