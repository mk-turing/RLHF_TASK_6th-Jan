package _12251

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestFeatureToggles(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	func isFeatureEnabled(feature string) bool {
		value, ok := os.LookupEnv(feature)
		if !ok {
		return false // Default to false if toggle is not defined in the config
	}
		enabled, err := strconv.ParseBool(value)
		if err != nil {
		t.Fatalf("Invalid value for %s toggle: %s", feature, value)
	}
		return enabled
	}

	t.Run("Disabled when not set", func(t *testing.T) {
		if isFeatureEnabled("NEW_FEATURE") {
			t.Error("NEW_FEATURE should be disabled when not set")
		}
	})

	t.Run("Enabled when set to true", func(t *testing.T) {
		os.Setenv("IMPROVED_ALGORITHM", "true")
		defer os.Unsetenv("IMPROVED_ALGORITHM")

		if !isFeatureEnabled("IMPROVED_ALGORITHM") {
			t.Error("IMPROVED_ALGORITHM should be enabled when set to true")
		}
	})

	t.Run("Disabled when set to false", func(t *testing.T) {
		os.Setenv("BUG_FIX", "false")
		defer os.Unsetenv("BUG_FIX")

		if isFeatureEnabled("BUG_FIX") {
			t.Error("BUG_FIX should be disabled when set to false")
		}
	})

	t.Run("Interdependent toggles", func(t *testing.T) {
		os.Setenv("NEW_UI", "true")
		defer os.Unsetenv("NEW_UI")

		if !isFeatureEnabled("NEW_UI") {
			t.Error("NEW_UI should be enabled when set to true")
		}

		if isFeatureEnabled("OLD_UI") {
			t.Error("OLD_UI should be disabled when NEW_UI is enabled")
		}
	})

	t.Run("Edge case: Invalid value", func(t *testing.T) {
		os.Setenv("FEATURE_X", "invalid")
		defer os.Unsetenv("FEATURE_X")

		if isFeatureEnabled("FEATURE_X") {
			t.Error("FEATURE_X should be disabled with an invalid value")
		}
	})
}