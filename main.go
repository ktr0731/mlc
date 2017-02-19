package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "invalid argument")
	}

	cmds := make([]*exec.Cmd, len(args[1:]))
	for i, arg := range args[1:] {
		cmdArgs := strings.Split(arg, " ")
		cmds[i] = exec.Command(cmdArgs[0], cmdArgs...)
	}

	exitChan := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("Start")
	go func() {
		<-sig
		fmt.Println("Interrupt")
		exitChan <- struct{}{}
	}()

	<-exitChan

	fmt.Println("Stop graceful")
}
