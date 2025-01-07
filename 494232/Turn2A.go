package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/spec"
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
		schemaV1Path string
		schemaV2Path string
	)
	flag.StringVar(&baseURL, "baseURL", "", "Base URL of the API")
	flag.StringVar(&versionParam, "versionParam", "", "Name of the version parameter in the query")
	flag.StringVar(&version1, "version1", "", "First version to compare")
	flag.StringVar(&version2, "version2", "", "Second version to compare")
	flag.StringVar(&queryV1, "queryV1", "", "Query parameters for version 1 as JSON string")
	flag.StringVar(&queryV2, "queryV2", "", "Query parameters for version 2 as JSON string")
	flag.StringVar(&schemaV1Path, "schemaV1", "schema_v1.json", "Path to schema file for version 1")
	flag.StringVar(&schemaV2Path, "schemaV2", "schema_v2.json", "Path to schema file for version 2")
	flag.Parse()

	if baseURL == "" || versionParam == "" || version1 == "" || version2 == "" || queryV1 == "" || queryV2 == "" {
		flag.Usage()
		return
	}

	// Load Swagger schemas
	schemaV1, err := loadSwaggerFile(schemaV1Path)
	if err != nil {
		fmt.Printf("Error loading schema v1: %v\n", err)
		return
	}

	schemaV2, err := loadSwaggerFile(schemaV2Path)
	if err != nil {
		fmt.Printf("Error loading schema v2: %v\n", err)
		return
	}

	paramsV1 := make(Parameters)
	err = json.Unmarshal([]byte(queryV1), &paramsV1)
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

	if err = validateQueryAgainstSchema(schemaV1, paramsV1); err != nil {
		fmt.Printf("Validation failed for version 1: %v\n", err)
		return
	}

	if err = validateQueryAgainstSchema(schemaV2, paramsV2); err != nil {
		fmt.Printf("Validation failed for version 2: %v\n", err)
		return
	}

	compResults := compareQueryParameters(paramsV1, paramsV2)
	fmt.Println("Comparison Results:")
	fmt.Println(compResults)
}

func loadSwaggerFile(path string) (*spec.Swagger, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var swagger spec.Swagger
	err = json.NewDecoder(file).Decode(&swagger)
	if err != nil {
		return nil, err
	}
	return &swagger, nil
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

func validateRequest(schema *spec.Swagger, path string, method string, params Parameters) error {
	// Resolve the path to the operation
	operation, found := getOperationForMethod(schema, path, method)
	if !found {
		return fmt.Errorf("operation for %s %s not found", method, path)
	}

	// Validate parameters
	for _, param := range operation.Parameters {
		if param.In == "query" { // We're dealing with query parameters
			// Check if the parameter exists in the request
			value, exists := params[param.Name]
			if !exists && param.Required {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}

			// Check type if needed (assuming string for simplicity, you may want to handle other types)
			if err := validateParamType(value, param.Type); err != nil {
				return fmt.Errorf("parameter '%s' validation failed: %v", param.Name, err)
			}
		}
	}

	// If all checks pass
	return nil
}

// Get operation for a specific HTTP method (GET, POST, etc.)
func getOperationForMethod(schema *spec.Swagger, path string, method string) (*spec.Operation, bool) {
	// Resolve the path
	paths := schema.Paths
	operation, exists := paths.Paths[path]
	if !exists {
		return nil, false
	}

	// Check for the corresponding method (GET, POST, etc.)
	switch method {
	case "GET":
		return operation.Get, true
	case "POST":
		return operation.Post, true
		// Add other HTTP methods as needed (PUT, DELETE, etc.)
	}

	return nil, false
}

// Validate parameter types (simplified for this example)
func validateParamType(value interface{}, paramType string) error {
	// Check if type matches (simplified, real-world cases would need more comprehensive checks)
	switch paramType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "integer":
		if _, ok := value.(int); !ok {
			return fmt.Errorf("expected integer, got %T", value)
		}
		// Add other types as necessary (boolean, float, etc.)
	}
	return nil
}

func validateQueryAgainstSchema(schema *spec.Swagger, params Parameters) error {
	filepath.Abs(filepath.Dir(os.Args[0]))
	// Resolve the relative path using the Swagger spec
	path := "/some-endpoint" // Replace with the actual path needed
	operation := schema.Paths.Paths[path].Get

	if operation == nil {
		return fmt.Errorf("operation not found for path %s", path)
	}

	req := &http.Request{}
	req.URL = &url.URL{Path: path, RawQuery: paramsToQueryString(params)}

	return validateRequest(schema, path, "GET", params)
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

func paramsToQueryString(params Parameters) string {
	queryPairs := []string{}
	for k, v := range params {
		queryPairs = append(queryPairs, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(queryPairs, "&")
}
