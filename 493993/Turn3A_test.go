package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSomeFunction(t *testing.T) {
	// Set up logger
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	// Hook up a test hook to capture logs
	hook := test.NewLocal(log)
	log.Hooks.Add(hook)

	// Call the function being tested
	err := someFunction(log)
	require.Error(t, err) // Assert that an error is returned

	// Assert that the correct error message is logged
	// Check that there are exactly 2 entries logged (if needed)
	assert.Equal(t, 4, len(hook.Entries)) // Check for two log entries

	// Ensure that the correct log levels and messages are captured
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)                                  // Ensure log level is Error
	assert.Equal(t, "An error occurred while processing the request", hook.Entries[0].Message) // Assert the error message

	// Additional assertion for the second log entry (if needed)
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[2].Level)           // Ensure second log entry is also Error level
	assert.Contains(t, hook.Entries[2].Message, "something went wrong") // Ensure the second message contains the error
}

func someFunction(log *logrus.Logger) error {
	// Simulate an error and log it
	err := fmt.Errorf("something went wrong")
	log.Error("An error occurred while processing the request")
	log.Error(err)
	return err
}
