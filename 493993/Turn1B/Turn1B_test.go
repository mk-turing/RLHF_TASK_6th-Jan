package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSomeFunction(t *testing.T) {
	// Hook up a logrus testing hook
	hook := test.NewLocal(logrus.StandardLogger())

	// Call the function you want to test
	//if err := someFunction(); err != nil {
	//	logrus.WithError(err).Error("An error occurred while processing the request")
	//}

	// Assert that the log message was written
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
	assert.Equal(t, "An error occurred while processing the request", hook.Entries[0].Message)
}
