package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Command is an interface for chat commands.
type Command interface {
	Execute(user string, args []string) string
}

// sendCommand represents the "/send" command.
type sendCommand struct{}

func (c sendCommand) Execute(user string, args []string) string {
	if len(args) == 0 {
		return "Invalid command. Usage: /send <message>"
	}
	message := strings.Join(args, " ")
	return fmt.Sprintf("%s: %s", user, message)
}

// helpCommand represents the "/help" command.
type helpCommand struct{}

func (c helpCommand) Execute(user string, args []string) string {
	return `Available commands:
/send <message>: Sends a message to the chat.
/help: Shows this help message.
/history: Shows the command history.`
}

// historyCommand represents the "/history" command.
type historyCommand struct {
	log *[]string
	mu  sync.RWMutex
}

func (c *historyCommand) Execute(user string, args []string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return strings.Join(*c.log, "\n")
}

// AppendToLog adds a command to the log.
func (c *historyCommand) AppendToLog(command string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	*c.log = append(*c.log, command)
}

// commands is a slice of registered commands.
var commands = []Command{
	sendCommand{},
	helpCommand{},
	&historyCommand{log: &[]string{}},
}

// user represents a user in the chat.
type user struct {
	name string
}

// chat represents the chat application.
type chat struct {
	users   map[string]*user
	history *historyCommand
}

// newChat creates a new chat application.
func newChat() *chat {
	historyCmd := commands[2].(*historyCommand)
	return &chat{
		users:   make(map[string]*user),
		history: historyCmd,
	}
}

// parseCommand parses the input message and returns the corresponding command and arguments.
func (c *chat) parseCommand(input string) (Command, []string) {
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
		case *historyCommand:
			if commandName == "/history" {
				return c, args
			}
		}
	}
	return nil, nil // Invalid command
}

// executeCommand executes the given command for the specified user and returns the response.
func (c *chat) executeCommand(user *user, command Command, args []string) string {
	response := command.Execute(user.name, args)
	// Log the command execution
	c.history.AppendToLog(fmt.Sprintf("%s: %s", user.name, input))
	return response
}

func main() {
	fmt.Println("Welcome to the chat application!")
	scanner := bufio.NewScanner(os.Stdin)
	chatApp := newChat()

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "/exit" {
			fmt.Println("Goodbye!")
			break
		}

		// TODO: Implement user login/logout system to manage users in the chat.
		currentUser := &user{name: "Anonymous"} // For simplicity, assume every user is "Anonymous"