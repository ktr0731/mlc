package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "invalid argument")
		return
	}

	wg := new(sync.WaitGroup)
	cmds := make([]*exec.Cmd, len(args[1:]))
	for i, arg := range args[1:] {
		wg.Add(2) // stdout and stderr
		cmdArgs := strings.Split(arg, " ")
		cmds[i] = exec.Command("sh", append([]string{"-c"}, strings.Join(cmdArgs, " "))...)
	}

	// TODO: go-shellwords を使う

	for _, cmd := range cmds {
		out, err := cmd.StdoutPipe()
		if err != nil {
			logging("cannot get stdout")
			return
		}
		errp, err := cmd.StderrPipe()
		if err != nil {
			logging("cannot get stderr")
			return
		}

		go scan(wg, out, os.Stdout)
		go scan(wg, errp, os.Stderr)

		if err := cmd.Start(); err != nil {
			logging(fmt.Sprintf("error while starting child processes: %s", err))
			// TODO: 起動してるプロセスを閉じる
			return
		}
	}

	exitChan := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	logging("Start")
	go func() {
		<-sig

		logging("Interrupt")
		logging("Send interrupt signal to children")

		for _, cmd := range cmds {
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				logging(fmt.Sprintf("error while killing child processes: %s", err))
				return
			}

			_, err := cmd.Process.Wait()
			if err != nil {
				logging(fmt.Sprintf("error while waiting child process: %s", err))
			}
		}

		logging("childlen processes killed")

		exitChan <- struct{}{}
	}()

	go func() {
		wg.Wait()
		logging("all processes terminated")
		exitChan <- struct{}{}
	}()

	<-exitChan

	logging("Stop")
}

func scan(wg *sync.WaitGroup, reader io.Reader, out io.Writer) {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		fmt.Fprintln(out, s.Text())
	}
	wg.Done()
}

func logging(message string) {
	if debug := os.Getenv("DEBUG"); debug != "" {
		log.Println(message)
	}
}
