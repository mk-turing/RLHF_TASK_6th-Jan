package _12251

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

type FeatureConfig struct {
	Feature1     bool                `json:"feature1"`
	Feature2     bool                `json:"feature2"`
	Feature3     bool                `json:"feature3"`
	Dependencies map[string][]string `json:"dependencies"`
}

func runTestsWithToggles(t *testing.T) {
	configFile := "turn3B.json"
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config FeatureConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// Validate dependencies
	if err := validateDependencies(&config); err != nil {
		t.Fatalf("Invalid dependencies: %v", err)
	}

	wg := sync.WaitGroup{}
	testCh := make(chan struct{}, 10)

	for _, toggle1 := range []bool{true, false} {
		for _, toggle2 := range []bool{true, false} {
			for _, toggle3 := range []bool{true, false} {
				wg.Add(1)
				go func(toggle1, toggle2, toggle3 bool) {
					defer wg.Done()
					c := config
					c.Feature1 = toggle1
					c.Feature2 = toggle2
					c.Feature3 = toggle3
					testCh <- struct{}{}
					defer func() { <-testCh }()
					runTest(t, &c)
				}(toggle1, toggle2, toggle3)
			}
		}
	}

	wg.Wait()
}

func runTest(t *testing.T, config *FeatureConfig) {
	t.Run(fmt.Sprintf("Feature1: %t, Feature2: %t, Feature3: %t", config.Feature1, config.Feature2, config.Feature3), func(t *testing.T) {
		if config.Feature1 {
			testFeature1(t)
		}
		if config.Feature2 {
			testFeature2(t)
		}
		if config.Feature3 {
			testFeature3(t)
		}
	})
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

func validateDependencies(config *FeatureConfig) error {
	for feature, dependencies := range config.Dependencies {
		for _, dep := range dependencies {
			switch dep {
			case "Feature1":
				if !config.Feature1 {
					return fmt.Errorf("%s requires Feature1 to be enabled", feature)
				}
			case "Feature2":
				if !config.Feature2 {
					return fmt.Errorf("%s requires Feature2 to be enabled", feature)
				}
			case "Feature3":
				if !config.Feature3 {
					return fmt.Errorf("%s requires Feature3 to be enabled", feature)
				}
			default:
				return fmt.Errorf("Invalid dependency %q for feature %q", dep, feature)
			}
		}
	}
	return nil
}

func TestFeatureToggles(t *testing.T) {
	runTestsWithToggles(t)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
