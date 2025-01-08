package main

import (
	"github.com/sirupsen/logrus"
)

const panicLevel logrus.Level = 5

func main() {
	// Register the custom panic level
	logrus.AddLevel(panicLevel, logrus.ErrorColor, "PANIC")

	// Set the log level to panic
	logrus.SetLevel(panicLevel)

	logrus.Debug("This is a debug message.")  // Will not be printed
	logrus.Info("This is an info message.")   // Will not be printed
	logrus.Warn("This is a warn message.")    // Will not be printed
	logrus.Error("This is an error message.") // Will not be printed
	logrus.Panic("This is a panic message.")  // Will be printed and the application will exit
}
