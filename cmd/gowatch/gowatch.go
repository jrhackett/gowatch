package main

import (
	"flag"
	"fmt"
	"gowatch/internal/watcher"
	"log"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	command := flag.String("command", "", "a command to run when any file inside of path is changed")
	path := flag.String("path", "./", "a path to watch recursively")

	flag.Parse()

	if *command == "" {
		color.Red("Missing command")
		return
	}

	watcher, err := watcher.NewFileWatcher(*path)
	if err != nil {
		log.Fatal(err)
	}

	go watcher.Run()

	for {
		<-watcher.Files
		commands := strings.Split(*command, " ")

		if len(commands) >= 1 {
			cmd := exec.Command(commands[0], commands[1:]...)
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
	}
}
