package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// truncateString truncates a string to maxLength characters
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// truncateValue recursively processes JSON values and truncates strings
func truncateValue(v interface{}, maxLength int) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range val {
			result[key] = truncateValue(value, maxLength)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = truncateValue(item, maxLength)
		}
		return result
	case string:
		return truncateString(val, maxLength)
	default:
		return v
	}
}

func main() {
	// Read the original JSON file
	data, err := ioutil.ReadFile("extracted_content/extracted_content.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Parse JSON
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Truncate all text values
	truncatedData := truncateValue(jsonData, 250)

	// Marshal back to JSON with pretty formatting
	outputData, err := json.MarshalIndent(truncatedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write to new file
	err = ioutil.WriteFile("extracted_content/extracted_content_short.json", outputData, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("Successfully created extracted_content_short.json with truncated text values (250 characters max)")
}
