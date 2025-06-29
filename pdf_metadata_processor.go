package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileMetadata represents the structure of each entry in the metadata JSON
type FileMetadata struct {
	FileName    string   `json:"file_name"`
	Title       string   `json:"Title"`
	Genre       string   `json:"Genre"`
	Tags        []string `json:"Tags"`
	Composer    string   `json:"Composer"`
	NewFilename string   `json:"new_filename"`
}

func main() {
	// Read metadata from JSON file
	metadata, err := readMetadata("extracted_content/metadata_clean.json")
	if err != nil {
		log.Fatalf("Error reading metadata: %v", err)
	}

	// Create output directory
	outputDir := "processed_files_with_metadata"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Check if exiftool is available
	exiftoolAvailable := checkExiftool()

	// Process each file
	processedCount := 0
	errorCount := 0

	for _, fileMeta := range metadata {
		if err := processFileWithMetadata(fileMeta, outputDir, exiftoolAvailable); err != nil {
			log.Printf("Error processing %s: %v", fileMeta.FileName, err)
			errorCount++
			continue
		}
		fmt.Printf("Processed with metadata: %s -> %s\n", fileMeta.FileName, fileMeta.NewFilename)
		processedCount++
	}

	fmt.Printf("\nProcessing complete!\n")
	fmt.Printf("Successfully processed: %d files\n", processedCount)
	fmt.Printf("Errors: %d files\n", errorCount)
	fmt.Printf("Files saved to: %s/\n", outputDir)

	if exiftoolAvailable {
		fmt.Printf("\n✅ PDF metadata was modified and should be visible in Mac Preview!\n")
	} else {
		fmt.Printf("\n⚠️  exiftool not found. Files were copied but metadata not embedded.\n")
		fmt.Printf("   Install exiftool with: brew install exiftool\n")
	}
}

func readMetadata(filename string) ([]FileMetadata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata []FileMetadata
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func checkExiftool() bool {
	_, err := exec.LookPath("exiftool")
	return err == nil
}

func processFileWithMetadata(fileMeta FileMetadata, outputDir string, useExiftool bool) error {
	// Source file path
	sourcePath := filepath.Join("extracted_content", fileMeta.FileName)

	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", sourcePath)
	}

	// Output file path
	outputPath := filepath.Join(outputDir, fileMeta.NewFilename)

	// Copy the file first
	if err := copyFile(sourcePath, outputPath); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	// If exiftool is available, modify the metadata
	if useExiftool {
		if err := modifyPdfMetadata(outputPath, fileMeta); err != nil {
			return fmt.Errorf("error modifying PDF metadata: %v", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func modifyPdfMetadata(filePath string, fileMeta FileMetadata) error {
	// Prepare exiftool command
	args := []string{
		"-overwrite_original",
		"-Title=" + fileMeta.Title,
		"-Author=" + fileMeta.Composer,
		"-Subject=" + fileMeta.Genre,
		"-Keywords=" + strings.Join(fileMeta.Tags, ", "),
		filePath,
	}

	cmd := exec.Command("exiftool", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
