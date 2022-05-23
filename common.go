package shell

type commonShell struct {
	command string
	args    string
	err     error
	isStop  bool
}

type CommonSheller interface {
	GetCommand() string
	GetArgs() string
	GetError() error
	HasStop() bool
	HasError() bool
}

func (c *commonShell) GetCommand() string {
	return c.command
}

func (c *commonShell) GetArgs() string {
	return c.args
}

func (c *commonShell) GetError() error {
	return c.err
}

func (c *commonShell) HasStop() bool {
	return c.isStop
}
func (c *commonShell) HasError() bool {
	return c.err != nil
}
