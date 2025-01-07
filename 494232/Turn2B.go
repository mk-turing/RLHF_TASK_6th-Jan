package main
import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)
type Parameters map[string]interface{}
// APIQuerySchema defines the schema for query parameters
type APIQuerySchema struct {
	RequiredParameters []string `json:"required_parameters"`
	OptionalParameters []string `json:"optional_parameters"`
	Parameters         map[string]string `json:"parameters"`
}
func main() {
	var (
		baseURL      string
		versionParam string
		version1     string
		version2     string
		queryV1      string
		queryV2      string
		schemaV1     string
		schemaV2     string
		errorLogPath  string
	)

	// Define the flags with help messages
	flag.StringVar(&baseURL, "baseURL", "", "Base URL of the API")
	flag.StringVar(&versionParam, "versionParam", "", "Name of the version parameter in the query")
	flag.StringVar(&version1, "version1", "", "First version to compare")
	flag.StringVar(&version2, "version2", "", "Second version to compare")
	flag.StringVar(&queryV1, "queryV1", "", "Query parameters for version 1 as JSON string")
	flag.StringVar(&queryV2, "queryV2", "", "Query parameters for version 2 as JSON string")
	flag.StringVar(&schemaV1, "schemaV1", "", "Path to JSON schema for version 1")
	flag.StringVar(&schemaV2, "schemaV2", "", "Path to JSON schema for version 2")
	flag.StringVar(&errorLogPath, "errorLogPath", "", "Path to generate error log file")
	flag.Parse()
	// Check for required flags
	if baseURL == "" || versionParam == "" || version1 == "" || version2 == "" || queryV1 == "" || queryV2 == "" {
		flag.Usage()
		return
	}
	// Define the error file
	errorFile, err := os.Create(errorLogPath)
	if err != nil {
		fmt.Println("Error creating error log file:", err)
		return
	}
	defer errorFile.Close()
	// Initialize csv writer for error log
	csvWriter := csv.NewWriter(errorFile)
	defer csvWriter.Flush()
	errorHeader := []string{"Timestamp", "Version", "Parameter", "Error"}
	if err := csvWriter.Write(errorHeader); err != nil {
		fmt.Println("Error writing error header:", err)
		return
	}
	// Read and validate the query schemas for version 1 and 2
	schema1 := parseSchema(schemaV1, version1, csvWriter)
	schema2 := parseSchema(schemaV2, version2, csvWriter)
	if schema1 == nil || schema2 == nil {
		return
	}
	// Parse query parameters for version 1 and 2
	paramsV1 := make(Parameters)
	if err := json.Unmarshal([]byte(queryV1), &paramsV1); err != nil {
		logError(csvWriter, version1, "Query Parameter Unmarshalling", err.Error())
		return
	}
	paramsV2 := make(Parameters)
	if err := json.Unmarshal([]byte(queryV2), &paramsV2); err != nil {
		logError(csvWriter, version2, "Query Parameter Unmarshalling", err.Error())
		return
	}
	if !checkBackwardCompatibility(fmt.Sprintf("%s?%s=%s", baseURL, versionParam, version1), paramsV1, paramsV2) {
		logError(csvWriter, "Both", "Backward Compatibility", "Version 2 not backward compatible with version 1")