package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

// SchemaVersion represents a unique version of a schema
type SchemaVersion string

const (
	SchemaVersion1 SchemaVersion = "SchemaVersion1"
	SchemaVersion2 SchemaVersion = "SchemaVersion2"
)

// SourceData represents data in the source system
type SourceData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TransformedData represents data after transformation
type TransformedData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	AgeGroup string `json:"age_group"`
	IsAdult  bool   `json:"is_adult"`
}

// Schema defines the structure of the data
type Schema struct {
	Version SchemaVersion `json:"version"`
	Fields  map[string]string `json:"fields"`
}

func main() {
	// Initialize schema registry
	schemaRegistry := NewRedisSchemaRegistry()

	// Initialize data pipeline
	pipeline := NewDataPipeline(schemaRegistry)

	// Start processing data
	pipeline.Start()
}

// DataPipeline represents the ETL pipeline
type DataPipeline struct {
	schemaRegistry SchemaRegistry
}

// NewDataPipeline creates a new DataPipeline
func NewDataPipeline(schemaRegistry SchemaRegistry) *DataPipeline {
	return &DataPipeline{
		schemaRegistry: schemaRegistry,
	}
}

// Start starts the data processing pipeline
func (p *DataPipeline) Start() {
	// Extract data from source
	sourceData := extractData()

	// Transform data based on latest schema
	transformedData := transformData(p.schemaRegistry, sourceData)

	// Load transformed data
	loadData(transformedData)
}

func extractData() []SourceData {
	return []SourceData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 18},
	}
}

func transformData(schemaRegistry SchemaRegistry, sourceData []SourceData) []TransformedData {
	var transformedData []TransformedData