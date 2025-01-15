package _12251

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

type Feature struct {
	Name         string   `json:"name"`
	State        bool     `json:"state"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type FeatureConfig struct {
	Features []Feature `json:"features"`
}

func runTestsWithToggles(t *testing.T) {
	configFile := "turn3A.json"
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config FeatureConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if !validateConfiguration(t, &config) {
		return
	}

	wg := sync.WaitGroup{}
	testCh := make(chan struct{}, 10) // Limit concurrent tests

	for _, feature := range config.Features {
		wg.Add(1)
		go func(feature Feature) {
			defer wg.Done()
			testCh <- struct{}{}        // Acquire a channel slot to limit concurrency
			defer func() { <-testCh }() // Release the slot after test completes

			runTest(t, &feature)
		}(feature)
	}

	wg.Wait()
}

func runTest(t *testing.T, feature *Feature) {
	t.Run(fmt.Sprintf("Feature: %s, State: %t", feature.Name, feature.State), func(t *testing.T) {
		switch feature.Name {
		case "feature1":
			testFeature1(t)
		case "feature2":
			testFeature2(t)
		case "feature3":
			testFeature3(t)
		default:
			t.Fatalf("Unknown feature: %s", feature.Name)
		}
	})
}

func validateConfiguration(t *testing.T, config *FeatureConfig) bool {
	for _, feature := range config.Features {
		for _, dep := range feature.Dependencies {
			found := false
			for _, f := range config.Features {
				if f.Name == dep {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Dependency '%s' not found for feature '%s'", dep, feature.Name)
				return false
			}

			if f := findFeatureByName(config, dep); f != nil && !f.State {
				t.Errorf("Dependency '%s' must be true for feature '%s'", dep, feature.Name)
				return false
			}
		}
	}

	return true
}

func findFeatureByName(config *FeatureConfig, name string) *Feature {
	for _, f := range config.Features {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func testFeature1(t *testing.T) {
	// Test logic for feature 1
}

func testFeature2(t *testing.T) {
	// Test logic for feature 2
}

func testFeature3(t *testing.T) {
	// Test logic for feature 3
}

func TestFeatureToggles(t *testing.T) {
	// Run the tests with feature toggles
	runTestsWithToggles(t)
}

func TestMain(m *testing.M) {
	// Run the main test suite
	os.Exit(m.Run())
}
