package _12231

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// UserRole represents the access level of a user.
type UserRole int

const (
	UserRoleUser UserRole = iota
	UserRoleAdmin
)

// Command is an interface for chat commands.
type Command interface {
	Execute(user string, args []string) string
	HasPermission(userRole UserRole) bool
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

func (c sendCommand) HasPermission(userRole UserRole) bool {
	return true // All users can send messages
}

// helpCommand represents the "help" command.
type helpCommand struct{}

func (c helpCommand) Execute(user string, args []string) string {
	return `Available commands:
/send <message>: Sends a message to the chat.
/help: Shows this help message.
/history: Shows the command history.`
}

func (c helpCommand) HasPermission(userRole UserRole) bool {
	return true // All users can access help
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
	return strings.Join(history, "\n")
}

func (c historyCommand) HasPermission(userRole UserRole) bool {
	return userRole == UserRoleAdmin // Only admins can access history
}

// CommandLog is a structure to hold a command log entry.
type CommandLog struct {
	User     string
	Command  string
	Timestamp time.Time
}

// User represents a user in the chat application.
type User struct {
	Name  string
	Role  UserRole
	mutex sync.Mutex
}

// NewUser creates a new user with the specified name and role.
func NewUser(name string, role UserRole) *User {
	return &User{Name: name, Role: role}
}

// GetName returns the user's name.
func (u *User) GetName() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.Name
}

// SetName sets the user's name.
func (u *User) SetName(name string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.Name = name
}

// GetRole returns the user's role.
func (u *User) GetRole() UserRole {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.Role
}

// SetRole sets the user's role.
func (u *User) SetRole(role UserRole) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.Role = role
}

// commands is a slice of registered commands.
var commands = []Command{
	sendCommand{},
	helpCommand{},
	historyCommand{},
}

// commandLog is a slice to log command executions.
var commandLog []CommandLog

// users is a map to store registered users.
var users = make(map[string]*User)

// parseCommand parses the input message and returns the corresponding command and arguments.
func parseCommand(input string) (Command, []string) {
	parts := strings.Split(input, " ")
	commandName := parts[0]
	args := parts[1:]

	for _, cmd := range commands {
		switch c := cmd.(type) {
		case sendCommand: