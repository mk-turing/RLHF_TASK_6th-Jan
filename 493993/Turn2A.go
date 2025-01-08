package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"strings"
)

func main() {
	// Set up logger
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Hook up a logrus testing hook
	hook := test.NewLocal(log)

	// Log various messages
	log.Debug("Debug message")
	log.Info("Info message")
	log.Warn("Warn message")
	log.Error("Error message: something went wrong")
	log.Fatal("Fatal error: application crashing")

	// Filter logs by log level (focus on Error and Fatal)
	filteredLogs := filterLogsByLevel(hook.Entries, logrus.ErrorLevel, logrus.FatalLevel)
	printFilteredLogs(filteredLogs)

	// Filter logs by specific keyword ("error")
	filteredLogsByKeyword := filterLogsByKeyword(hook.Entries, "error")
	printFilteredLogs(filteredLogsByKeyword)
}

func filterLogsByLevel(entries []logrus.Entry, levels ...logrus.Level) []logrus.Entry {
	filtered := make([]logrus.Entry, 0)
	for _, entry := range entries {
		if containsLevel(levels, entry.Level) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func containsLevel(levels []logrus.Level, level logrus.Level) bool {
	for _, l := range levels {
		if l == level {
			return true
		}
	}
	return false
}

func filterLogsByKeyword(entries []logrus.Entry, keyword string) []logrus.Entry {
	filtered := make([]logrus.Entry, 0)
	keywordLower := strings.ToLower(keyword)
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Message), keywordLower) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func printFilteredLogs(entries []logrus.Entry) {
	for _, entry := range entries {
		// Use the Formatter to format the entry and print it
		fmt.Println(entry.String())
	}
}
