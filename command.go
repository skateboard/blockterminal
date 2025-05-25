package terminal

type Command interface {
	Name() string
	Description() string
	Args() []string
	Execute(args []string) error
}

type BasicCommand struct {
	commandName string
	description string
	args        []string
}

func newBasicCommand(commandName string, description string, args []string) BasicCommand {
	return BasicCommand{commandName: commandName, description: description, args: args}
}

func (c *BasicCommand) Name() string {
	return c.commandName
}

func (c *BasicCommand) Description() string {
	return c.description
}

func (c *BasicCommand) Args() []string {
	return c.args
}
