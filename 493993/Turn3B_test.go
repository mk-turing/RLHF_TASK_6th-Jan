package main

import (
	"errors"
	"github.com/sirupsen/logrus/hooks/test"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func someFunction() error {
	// Your application code here
	// ...

	return errors.New("something went wrong")
}

func TestSomeFunction(t *testing.T) {
	// Hook up a logrus testing hook
	hook := test.NewLocal(logrus.StandardLogger())

	// Call the function you want to test
	err := someFunction()
	assert.Error(t, err)

	// Assert that the log message was written
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
	assert.Equal(t, "something went wrong", hook.Entries[0].Message)
}
