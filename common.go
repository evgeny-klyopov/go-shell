package shell

type commonShell struct {
	command string
	args    string
	success bool
	err     error
	isStop  bool
}

type CommonSheller interface {
	GetCommand() string
	GetArgs() string
	GetSuccess() bool
	GetError() error
	HasStop() bool
}

func (c *commonShell) GetCommand() string {
	return c.command
}

func (c *commonShell) GetArgs() string {
	return c.args
}

func (c *commonShell) GetSuccess() bool {
	return c.success
}

func (c *commonShell) GetError() error {
	return c.err
}

func (c *commonShell) HasStop() bool {
	return c.isStop
}
