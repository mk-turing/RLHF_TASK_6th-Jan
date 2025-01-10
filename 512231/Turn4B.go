package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Command is an interface for chat commands.
type Command interface {
	Execute(user string, args []string) string
	GetPermissions() []string
}

// Role represents the user's role in the chat application.
type Role string

const (
	Guest  Role = "guest"
	Member Role = "member"
	Admin  Role = "admin"
)

// user contains the user's information.
type user struct {
	Name     string
	Role     Role
	commands []string
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

func (c sendCommand) GetPermissions() []string {
	return []string{"guest"}
}

// helpCommand represents the "help" command.
type helpCommand struct{}

func (c helpCommand) Execute(user string, args []string) string {
	return `Available commands:
/send <message>: Sends a message to the chat.
/help: Shows this help message.
/history: Shows the command history.
/bans: Lists banned users.`
}

func (c helpCommand) GetPermissions() []string {
	return []string{"guest"}
}

// historyCommand represents the "history" command.
type historyCommand struct{}

func (c historyCommand) Execute(user string, args []string) string {
	var history []string
	for _, logEntry := range commandLog {
		if logEntry.User == user {
			history = append(history, logEntry.Command)
		}
	}
	if len(history) == 0 {
		return "No command history found."
	}
	return strings.Join(history, "\n")
}

func (c historyCommand) GetPermissions() []string {
	return []string{"member", "admin"}
}

// bansCommand represents the "bans" command.
type bansCommand struct{}

func (c bansCommand) Execute(user string, args []string) string {
	if len(bannedUsers) == 0 {
		return "No banned users."
	}
	var bans []string
	for _, bannedUser := range bannedUsers {
		bans = append(bans, fmt.Sprintf("%s - Banned since %s", bannedUser.Name, bannedUser.BanTime.String()))
	}
	return strings.Join(bans, "\n")
}

func (c bansCommand) GetPermissions() []string {
	return []string{"admin"}
}

// CommandLog is a structure to hold a command log entry.
type CommandLog struct {
	User      string
	Command   string
	Timestamp time.Time
}

// bannedUser is a structure to hold information of a banned user.
type bannedUser struct {
	Name    string
	BanTime time.Time
}

var bannedUsers []bannedUser

// commands is a slice of registered commands.
var commands = []Command{
	sendCommand{},
	helpCommand{},
	historyCommand{},
	bansCommand{},
}

// commandLog is a slice to log command executions.
var commandLog []CommandLog

var commandLogMutex sync.Mutex

// userInputHandler handles user input and executes commands concurrently.
func userInputHandler(user user, inputChan <-chan string) {
	for input := range inputChan {
		command, args := parseCommand(input)
		if command == nil {
			fmt.Println("Invalid command. Type /help for help.")
		} else {
			response := executeCommand(user, command, args)
			fmt.Println(response)
		}
	}
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
		case historyCommand:
			if commandName == "/history" {
				return c, args
			}
		case bansCommand:
			if commandName == "/bans" {
				return c, args
			}
		}
	}
	return nil, nil // Invalid command
}

// executeCommand executes the given command for a user and returns the response.
func executeCommand(user user, command Command, args []string) string {
	response := command.Execute(user.Name, args)
	logCommand(user, command, args)
	return response
}

// logCommand logs the command execution.
func logCommand(user user, command Command, args []string) {
	commandLogMutex.Lock()
	defer commandLogMutex.Unlock()
	commandLog = append(commandLog, CommandLog{
		User:      user.Name,
		Command:   fmt.Sprintf("%s %v", command, args),
		Timestamp: time.Now(),
	})
}

func main() {
	fmt.Println("Welcome to the chat application!")
	users := []user{
		{"Alice", Member, nil},
		{"Bob", Admin, nil},
		{"Charlie", Guest, nil},
	}

	inputChans := make([]chan string, len(users))
	for i := range users {
		inputChans[i] = make(chan string)
		go userInputHandler(users[i], inputChans[i])
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "/exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Broadcast the input to all user input channels
		for _, inputChan := range inputChans {
			inputChan <- input
		}
	}

	// Close all input channels to terminate the user input handlers
	for _, inputChan := range inputChans {
		close(inputChan)
	}
}
