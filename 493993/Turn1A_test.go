package main

import (
	"bytes"
	"log/capture"
	"strings"
	"testing"
)

func TestMyFunction(t *testing.T) {
	var buf bytes.Buffer
	defer capture.New(&buf).Close()

	// Call your function here
	//MyFunction()

	// Assert on log output
	if !strings.Contains(buf.String(), "Expected log message") {
		t.Errorf("Expected log message not found in output")
	}
}
