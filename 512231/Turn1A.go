package _12231

import "fmt"

type Command interface {
	Execute() error
}

var commandQueue []Command

type HelpCommand struct{}

func (cmd HelpCommand) Execute() error {
	fmt.Println("Available commands: help, join, leave, message")
	return nil
}

type JoinCommand struct {
	room string
}

func (cmd JoinCommand) Execute() error {
	fmt.Printf("Joined room: %s\n", cmd.room)
	return nil
}

// Add more commands as needed

func executeCommands() {
	for len(commandQueue) > 0 {
		command := commandQueue[0]
		commandQueue = commandQueue[1:] // Remove the executed command

		err := command.Execute()
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
