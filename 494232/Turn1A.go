package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

type Parameters map[string]interface{}

func main() {
	var (
		baseURL      string
		versionParam string
		version1     string
		version2     string
		queryV1      string
		queryV2      string
	)
	flag.StringVar(&baseURL, "baseURL", "", "Base URL of the API")
	flag.StringVar(&versionParam, "versionParam", "", "Name of the version parameter in the query")
	flag.StringVar(&version1, "version1", "", "First version to compare")
	flag.StringVar(&version2, "version2", "", "Second version to compare")
	flag.StringVar(&queryV1, "queryV1", "", "Query parameters for version 1 as JSON string")
	flag.StringVar(&queryV2, "queryV2", "", "Query parameters for version 2 as JSON string")
	flag.Parse()

	if baseURL == "" || versionParam == "" || version1 == "" || version2 == "" || queryV1 == "" || queryV2 == "" {
		flag.Usage()
		return
	}

	paramsV1 := make(Parameters)
	err := json.Unmarshal([]byte(queryV1), &paramsV1)
	if err != nil {
		fmt.Printf("Error parsing query1: %v\n", err)
		return
	}

	paramsV2 := make(Parameters)
	err = json.Unmarshal([]byte(queryV2), &paramsV2)
	if err != nil {
		fmt.Printf("Error parsing query2: %v\n", err)
		return
	}

	if !checkBackwardCompatibility(fmt.Sprintf("%s?%s=%s", baseURL, versionParam, version1), paramsV1, paramsV2) {
		fmt.Println("Backward compatibility not maintained!")
		return
	}

	compResults := compareQueryParameters(paramsV1, paramsV2)
	fmt.Println("Comparison Results:")
	fmt.Println(compResults)
}

func checkBackwardCompatibility(url string, v1Query Parameters, v2Query Parameters) bool {
	for k, v := range v1Query {
		_, exists := v2Query[k]
		if !exists || v2Query[k].(string) != v.(string) {
			return false
		}
	}
	return true
}

func compareQueryParameters(v1Query Parameters, v2Query Parameters) string {
	comparisonResults := []string{}

	for k, _ := range v1Query {
		if _, exists := v2Query[k]; !exists {
			comparisonResults = append(comparisonResults, "Version 2 misses parameter "+k)
		}
	}

	for k, _ := range v2Query {
		if _, exists := v1Query[k]; !exists {
			comparisonResults = append(comparisonResults, "Version 1 misses parameter "+k)
		}
	}

	return strings.Join(comparisonResults, "\n")
}
