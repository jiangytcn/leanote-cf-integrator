package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	. "github.com/yacloud-io/leanote-cf-integrator/mci"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	configFilePath := path.Join(wd, "conf", "app.conf")
	config, err := ExtractConfig(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = CFConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	PushConfig(config, configFilePath)
	var integratorExec *exec.Cmd
	if len(os.Args) > 1 {
		args := os.Args[1:]
		integratorExec = exec.Command(path.Join(wd, "bin", "platform"), args...)
	} else {
		integratorExec = exec.Command(path.Join(wd, "bin", "platform"), "-importPath github.com/leanote/leanote")
	}

	integratorExec.Stdout = os.Stdout
	integratorExec.Stderr = os.Stderr
	err = integratorExec.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
