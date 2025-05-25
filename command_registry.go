package terminal

import "fmt"

type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]Command)}
}

func (r *CommandRegistry) RegisterCommand(command Command) {
	r.commands[command.Name()] = command
}

func (r *CommandRegistry) GetCommand(name string) Command {
	return r.commands[name]
}

func (r *CommandRegistry) GetCommands() []Command {
	commands := make([]Command, 0, len(r.commands))
	for _, command := range r.commands {
		commands = append(commands, command)
	}

	return commands
}

func (r *CommandRegistry) RegisterCommands(commands []Command) {
	for _, command := range commands {
		r.RegisterCommand(command)
	}
}

func (r *CommandRegistry) RunCommand(name string, args []string) error {
	command := r.GetCommand(name)
	if command == nil {
		return fmt.Errorf("command %s not found", name)
	}

	return command.Execute(args)
}
