package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	watcher, err := NewFileWatcher("./")
	if err != nil {
		log.Fatal(err)
	}

	go watcher.run()

	for {
		<-watcher.Files
		goTest()
	}
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
