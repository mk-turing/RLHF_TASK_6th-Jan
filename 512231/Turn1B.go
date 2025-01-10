package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Command is an interface for chat commands.
type Command interface {
	Execute(user string, args []string) string
}

// sendCommand represents the "send" command.
type sendCommand struct{}

func (c sendCommand) Execute(user string, args []string) string {
	if len(args) == 0 {
		return "Invalid command. Usage: /send <message>"
	}
	message := strings.Join(args, " ")
	return fmt.Sprintf("%s: %s", user, message)
}

// helpCommand represents the "help" command.
type helpCommand struct{}

func (c helpCommand) Execute(user string, args []string) string {
	return `Available commands:
/send <message>: Sends a message to the chat.
/help: Shows this help message.`
}

// commands is a slice of registered commands.
var commands = []Command{
	sendCommand{},
	helpCommand{},
}

// parseCommand parses the input message and returns the corresponding command and arguments.
func parseCommand(input string) (Command, []string) {
	parts := strings.Split(input, " ")
	commandName := parts[0]
	args := parts[1:]

	for _, cmd := range commands {
		switch c := cmd.(type) {
		case sendCommand:
			if commandName == "/send" {
				return c, args
			}
		case helpCommand:
			if commandName == "/help" {
				return c, args
			}
		}
	}
	return nil, nil // Invalid command
}

// executeCommand executes the given command and returns the response.
func executeCommand(user string, command Command, args []string) string {
	return command.Execute(user, args)
}

func main() {
	fmt.Println("Welcome to the chat application!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "/exit" {
			fmt.Println("Goodbye!")
			break
		}

		command, args := parseCommand(input)
		if command == nil {
			fmt.Println("Invalid command. Type /help for help.")
		} else {
			user := "You"
			response := executeCommand(user, command, args)
			fmt.Println(response)
		}
	}
}
