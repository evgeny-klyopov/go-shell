package shell

import (
	"os/exec"
	"strings"
	"sync"
)

type shell struct {
	commonShell
	stop    chan bool
	readers map[string]reader
}

type Sheller interface {
	CommonSheller
	Stop()
	Run()
	HasChannel(typeChannel string) bool
	GetChannel(typeChannel string) chan string
	exec()
	close()
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
	go s.exec()
}

func (s *shell) Stop() {
	s.isStop = true
	s.stop <- s.isStop
}

func (s *shell) HasChannel(typeChannel string) bool {
	return s.readers[typeChannel].getEnable()
}

func (s *shell) GetChannel(typeChannel string) chan string {
	return s.readers[typeChannel].getChannel()
}

func (s *shell) exec() {
	var err error

	cmd := exec.Command(s.command, strings.Fields(s.args)...)

	for _, r := range s.readers {
		if r.getEnable() == true {
			err = r.setPipe(cmd)
			if err != nil {
				s.err = err
				s.Stop()
			}
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
			for _, r := range s.readers {
				if r.getEnable() == true {
					go r.startRead()
				}
			}
			wg.Done()
		}()

		wg.Wait()

		done <- cmd.Wait()
	}()
	select {
	case <-s.stop:
		s.close()
	case err = <-done:
		if err != nil {
			s.err = err
		}
		s.close()
		s.success = cmd.ProcessState.Success()
	}
}

func (s *shell) close() {
	for _, r := range s.readers {
		if r.getEnable() == true {
			r.closeChannel()
		}
	}
}
