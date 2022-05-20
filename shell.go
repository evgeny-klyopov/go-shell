package main

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"
)

type shell struct {
	command     string
	stop        chan bool
	output      chan string
	outputError chan string
	args        string
	success     bool
	err         error
	isStop      bool
	isOutput    bool
}

type Sheller interface {
	Stop()
	Run()
	GetOutput() chan string
	GetOutputFlag() bool
	GetCommand() string
	GetArgs() string
	GetSuccess() bool
	GetError() error
	GetOutputError() chan string

	exec()
	reader(r io.Reader, channel chan string)
}

func New(command string, args string, isOutput bool) Sheller {
	return &shell{
		command:  command,
		isOutput: isOutput,
		args:     args,
		stop:     make(chan bool),

		output:      make(chan string),
		outputError: make(chan string),
	}
}

func (s *shell) GetOutputFlag() bool {
	return s.isOutput
}
func (s *shell) GetCommand() string {
	return s.command
}

func (s *shell) GetArgs() string {
	return s.args
}

func (s *shell) GetSuccess() bool {
	return s.success
}

func (s *shell) GetError() error {
	return s.err
}

func (s *shell) GetOutputError() chan string {
	return s.outputError
}

func (s *shell) GetOutput() chan string {
	return s.output
}

func (s *shell) Run() {
	go s.exec()
}

func (s *shell) Stop() {
	s.isStop = true
	s.stop <- s.isStop
}

func (s *shell) exec() {
	cmd := exec.Command(s.command, strings.Fields(s.args)...)

	var stdOut, stdErr io.Reader
	var err error

	if s.isOutput == true {
		stdOut, err = cmd.StdoutPipe()
		if err != nil {
			s.err = err
			s.Stop()
		}

		stdErr, err = cmd.StderrPipe()
		if err != nil {
			s.err = err
			s.Stop()
		}
	}

	if err = cmd.Start(); err != nil {
		s.err = err
		s.Stop()
	}

	done := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			if s.isOutput == true {
				go s.reader(stdOut, s.output)
				go s.reader(stdErr, s.outputError)
			}
			wg.Done()
		}()

		wg.Wait()

		done <- cmd.Wait()
	}()
	select {
	case <-s.stop:
		close(s.output)
		close(s.outputError)
	case err = <-done:
		if err != nil {
			s.err = err
		}
		s.success = cmd.ProcessState.Success()
	}
}

func (s *shell) reader(r io.Reader, channel chan string) {
	reader := bufio.NewReader(r)
	for {
		str, err := reader.ReadString('\n')

		if err != nil {
			close(channel)
			break
		}
		channel <- strings.TrimSuffix(str, "\n")
	}
}
