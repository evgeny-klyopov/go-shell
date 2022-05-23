package shell

import (
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type shell struct {
	commonShell
	stop    chan bool
	readers map[string]reader
	cmd     *exec.Cmd
}

type Sheller interface {
	CommonSheller
	Stop()
	Run()
	HasChannel(typeChannel string) bool
	GetChannel(typeChannel string) chan string
	GetSuccess() bool
	exec()
}

func New(command string, args string, makeStdOut bool, makeStdErr bool) Sheller {
	return &shell{
		stop: make(chan bool),
		commonShell: commonShell{
			command: command,
			args:    args,
		},
		readers: map[string]reader{
			OutTypeChannel: newReader(OutTypeChannel, makeStdOut),
			ErrTypeChannel: newReader(ErrTypeChannel, makeStdErr),
		},
	}
}

func (s *shell) Run() {
	s.cmd = exec.Command(s.command, strings.Fields(s.args)...)
	go s.exec()
}

func (s *shell) Stop() {
	s.isStop = true
	s.stop <- s.isStop
}

func (s *shell) GetSuccess() bool {
	var success bool

	if s.cmd.ProcessState != nil && (s.HasStop() || s.HasError()) {
		success = s.cmd.ProcessState.Success()
	}

	return success
}

func (s *shell) HasChannel(typeChannel string) bool {
	return s.readers[typeChannel].getEnable()
}

func (s *shell) GetChannel(typeChannel string) chan string {
	return s.readers[typeChannel].getChannel()
}

func (s *shell) exec() {
	var err error

	for _, r := range s.readers {
		if r.getEnable() == true {
			err = r.setPipe(s.cmd)
			if err != nil {
				s.err = err
				s.Stop()
			}
		}
	}

	if err = s.cmd.Start(); err != nil {
		s.err = err
		s.Stop()
	}

	isError := make(chan error, 1)

	for _, r := range s.readers {
		if r.getEnable() == true {
			go r.startRead()
		}
	}

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
		}()

		wg.Wait()
		isError <- s.cmd.Wait()
	}()

	select {
	case <-s.stop:
		s.err = s.cmd.Process.Signal(syscall.SIGTERM)
		s.err = s.cmd.Process.Signal(syscall.SIGINT)
	case err = <-isError:
		if err != nil {
			s.err = err
		}
	}
}
