package shell

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

const OutTypeChannel = "out"
const ErrTypeChannel = "err"

type reader interface {
	getEnable() bool
	getType() string
	getChannel() chan string
	setPipe(cmd *exec.Cmd) error
	startRead()
	closeChannel()
}

type read struct {
	typeChannel string
	channel     chan string
	isEnable    bool
	pipe        io.Reader
}

func newReader(typeChannel string, isEnable bool) reader {
	return &read{
		typeChannel: typeChannel,
		isEnable:    isEnable,
		channel:     make(chan string),
	}
}

func (r *read) setPipe(cmd *exec.Cmd) error {
	var err error
	if r.typeChannel == OutTypeChannel {
		r.pipe, err = cmd.StdoutPipe()
	} else {
		r.pipe, err = cmd.StderrPipe()
	}

	return err
}

func (r *read) startRead() {
	pipe := bufio.NewReader(r.pipe)
	for {
		str, err := pipe.ReadString('\n')

		if err != nil {
			r.closeChannel()
			break
		}
		r.channel <- strings.TrimSuffix(str, "\n")
	}
}

func (r *read) closeChannel() {
	close(r.channel)
}

func (r *read) getEnable() bool {
	return r.isEnable
}

func (r *read) getType() string {
	return r.typeChannel
}

func (r *read) getChannel() chan string {
	return r.channel
}
