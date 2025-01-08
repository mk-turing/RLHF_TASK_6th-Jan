package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// Set the log level to debug
	logrus.SetLevel(logrus.DebugLevel)

	// Set the output of logs to stdout
	logrus.SetOutput(os.Stdout)

	// Enable log formatting
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logrus.Debug("This is a debug message")
	logrus.Info("This is an info message")
	logrus.Warn("This is a warn message")
	logrus.Error("This is an error message")
	logrus.Fatal("This is a fatal message")
}
