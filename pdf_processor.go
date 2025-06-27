package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type ExtractedContent struct {
	FileName    string    `json:"file_name"`
	ExtractedAt time.Time `json:"extracted_at"`
	Text        string    `json:"text"`
	PageCount   int       `json:"page_count"`
	FileSize    int64     `json:"file_size_bytes"`
}

type ProcessingResult struct {
	SourcePDF      string             `json:"source_pdf"`
	SplitFiles     []ExtractedContent `json:"split_files"`
	ProcessingTime time.Duration      `json:"processing_time_seconds"`
	TotalFiles     int                `json:"total_files_processed"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run pdf_processor.go <input_directory>")
		fmt.Println("Example: go run pdf_processor.go v4")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := "extracted_content"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Get all PDF files
	pdfFiles, err := filepath.Glob(filepath.Join(inputDir, "*.pdf"))
	if err != nil {
		log.Fatal("Failed to find PDF files:", err)
	}

	if len(pdfFiles) == 0 {
		log.Fatal("No PDF files found in", inputDir)
	}

	var allResults []ProcessingResult

	for _, pdfFile := range pdfFiles {
		fmt.Printf("Processing: %s\n", filepath.Base(pdfFile))

		startTime := time.Now()
		result := processPDF(pdfFile, outputDir)
		result.ProcessingTime = time.Since(startTime)

		allResults = append(allResults, result)
		fmt.Printf("Completed: %s (%d files created, %v)\n",
			filepath.Base(pdfFile), len(result.SplitFiles), result.ProcessingTime)
	}

	// Write results to JSON file
	outputFile := filepath.Join(outputDir, "extracted_content.json")
	jsonData, err := json.MarshalIndent(allResults, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON:", err)
	}

	if err := ioutil.WriteFile(outputFile, jsonData, 0644); err != nil {
		log.Fatal("Failed to write JSON file:", err)
	}

	fmt.Printf("\nProcessing complete! Results saved to: %s\n", outputFile)
	fmt.Printf("Total PDFs processed: %d\n", len(allResults))

	totalFiles := 0
	for _, result := range allResults {
		totalFiles += len(result.SplitFiles)
	}
	fmt.Printf("Total files extracted: %d\n", totalFiles)
}

func processPDF(pdfPath, outputDir string) ProcessingResult {
	result := ProcessingResult{
		SourcePDF: filepath.Base(pdfPath),
	}

	// Step 1: Find pages with "SPLITME"
	splitPages := findSplitPages(pdfPath)

	if len(splitPages) == 0 {
		fmt.Printf("  No 'SPLITME' pages found in %s\n", filepath.Base(pdfPath))
		return result
	}

	// Step 2: Split PDF into individual files
	splitFiles := splitPDFOnPages(pdfPath, outputDir, splitPages)

	// Step 3: Extract text from each split file
	for _, splitFile := range splitFiles {
		content := extractTextFromPDF(splitFile)
		result.SplitFiles = append(result.SplitFiles, content)
	}

	result.TotalFiles = len(result.SplitFiles)
	return result
}

func findSplitPages(pdfPath string) []int {
	// Use pdftotext to extract text and find pages with "SPLITME"
	cmd := exec.Command("pdftotext", "-f", "1", "-l", "999", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Warning: Failed to extract text from %s: %v", pdfPath, err)
		return nil
	}

	// Split by pages and find those containing "SPLITME"
	text := string(output)
	pages := strings.Split(text, "\f") // Form feed character separates pages

	var splitPages []int
	for i, page := range pages {
		if strings.Contains(strings.ToUpper(page), "SPLITME") {
			splitPages = append(splitPages, i+1) // Page numbers are 1-indexed
		}
	}

	return splitPages
}

func splitPDFOnPages(pdfPath, outputDir string, splitPages []int) []string {
	baseName := strings.TrimSuffix(filepath.Base(pdfPath), ".pdf")
	var splitFiles []string

	// Get total page count
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Warning: Failed to get PDF info for %s: %v", pdfPath, err)
		return nil
	}

	// Parse page count
	re := regexp.MustCompile(`Pages:\s+(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		log.Printf("Warning: Could not determine page count for %s", pdfPath)
		return nil
	}

	totalPages := 0
	fmt.Sscanf(matches[1], "%d", &totalPages)

	// Create ranges for splitting
	ranges := createPageRanges(splitPages, totalPages)

	// Split PDF using pdftk
	for i, pageRange := range ranges {
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_part_%03d.pdf", baseName, i+1))

		cmd := exec.Command("pdftk", pdfPath, "cat", pageRange, "output", outputFile)
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to split PDF %s range %s: %v", pdfPath, pageRange, err)
			continue
		}

		splitFiles = append(splitFiles, outputFile)
	}

	return splitFiles
}

func createPageRanges(splitPages []int, totalPages int) []string {
	var ranges []string
	start := 1

	for _, splitPage := range splitPages {
		if splitPage > start {
			ranges = append(ranges, fmt.Sprintf("%d-%d", start, splitPage-1))
		}
		start = splitPage + 1
	}

	// Add final range if there are pages after the last split
	if start <= totalPages {
		ranges = append(ranges, fmt.Sprintf("%d-%d", start, totalPages))
	}

	return ranges
}

func extractTextFromPDF(pdfPath string) ExtractedContent {
	content := ExtractedContent{
		FileName:    filepath.Base(pdfPath),
		ExtractedAt: time.Now(),
	}

	// Get file size
	if fileInfo, err := os.Stat(pdfPath); err == nil {
		content.FileSize = fileInfo.Size()
	}

	// First try to extract text directly (for text-based PDFs)
	cmd := exec.Command("pdftotext", pdfPath, "-")
	output, err := cmd.Output()

	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		content.Text = strings.TrimSpace(string(output))
		content.PageCount = countPages(pdfPath)
		return content
	}

	// If direct extraction fails or produces no text, use OCR
	fmt.Printf("  Using OCR for: %s\n", filepath.Base(pdfPath))
	content.Text = extractTextWithOCR(pdfPath)
	content.PageCount = countPages(pdfPath)

	return content
}

func extractTextWithOCR(pdfPath string) string {
	// Convert PDF to images first
	tempDir, err := ioutil.TempDir("", "pdf_ocr")
	if err != nil {
		log.Printf("Warning: Failed to create temp directory: %v", err)
		return ""
	}
	defer os.RemoveAll(tempDir)

	// Convert PDF to images using pdftoppm
	cmd := exec.Command("pdftoppm", "-png", "-r", "300", pdfPath, filepath.Join(tempDir, "page"))
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to convert PDF to images: %v", err)
		return ""
	}

	// Get all generated image files
	imageFiles, err := filepath.Glob(filepath.Join(tempDir, "page-*.png"))
	if err != nil {
		log.Printf("Warning: Failed to find image files: %v", err)
		return ""
	}

	var allText []string
	for _, imageFile := range imageFiles {
		// Use tesseract for OCR
		cmd := exec.Command("tesseract", imageFile, "stdout", "-l", "eng")
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Warning: OCR failed for %s: %v", imageFile, err)
			continue
		}

		text := strings.TrimSpace(string(output))
		if text != "" {
			allText = append(allText, text)
		}
	}

	return strings.Join(allText, "\n\n")
}

func countPages(pdfPath string) int {
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	re := regexp.MustCompile(`Pages:\s+(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return 0
	}

	var pageCount int
	fmt.Sscanf(matches[1], "%d", &pageCount)
	return pageCount
}
