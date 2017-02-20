package main

import (
	"bufio"
	"fmt"
	"io"
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
			// log.Println("cannot get stdout")
			return
		}
		errp, err := cmd.StderrPipe()
		if err != nil {
			// log.Println("cannot get stderr")
			return
		}

		go scan(wg, out, os.Stdout)
		go scan(wg, errp, os.Stderr)

		if err := cmd.Start(); err != nil {
			// log.Printf("error while starting child processes: %s", err)
			// TODO: 起動してるプロセスを閉じる
			return
		}
		// log.Println("PID:", cmd.Process.Pid)
		cmd.Stdout = os.Stdout
	}

	exitChan := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// log.Println("Start")
	go func() {
		<-sig
		// log.Println("Interrupt")
		// log.Println("Send signal to children")

		for _, cmd := range cmds {
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				// log.Printf("error while killing child processes: %s", err)
				return
			}

			_, err := cmd.Process.Wait()
			if err != nil {
				// log.Printf("error while waiting child process: %s", err)
			}

			// log.Println("exit on ", stat.String())
		}

		// log.Println("process killed")

		exitChan <- struct{}{}
	}()

	go func() {
		wg.Wait()
		// log.Println("processes finished")
		exitChan <- struct{}{}
	}()

	<-exitChan

	// log.Println("Stop graceful")
}

func scan(wg *sync.WaitGroup, reader io.Reader, out io.Writer) {
	defer func() {
		// log.Println("close scan")
	}()
	// log.Println("create scanner")
	s := bufio.NewScanner(reader)
	for s.Scan() {
		fmt.Fprintln(out, s.Text())
	}
	wg.Done()
}
