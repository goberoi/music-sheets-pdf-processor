# Music Sheets PDF Processor

A Go-based tool for processing guitar lesson PDFs, splitting them on "SPLITME" pages, and extracting text content using OCR.

## Features

- **PDF Splitting**: Automatically splits PDFs on pages containing "SPLITME"
- **Text Extraction**: Uses both direct text extraction and OCR (Tesseract) for maximum compatibility
- **JSON Output**: Generates structured JSON with metadata for each extracted file
- **Batch Processing**: Processes multiple PDFs in a directory
- **Cross-platform**: Works on macOS, Linux, and Windows

## Prerequisites

### Required Tools (install via Homebrew on macOS)

```bash
# PDF manipulation tools
brew install pdftk-java poppler

# OCR engine
brew install tesseract tesseract-lang
```

### Go

Make sure you have Go installed (version 1.16 or later recommended).

## Installation

1. Clone this repository:
```bash
git clone <your-repo-url>
cd music_sheets
```

2. Run the processor:
```bash
go run pdf_processor.go v4
```

## Usage

```bash
go run pdf_processor.go <input_directory>
```

### Example
```bash
go run pdf_processor.go v4
```

This will:
1. Process all PDF files in the `v4` directory
2. Find pages containing "SPLITME" 
3. Split PDFs at those pages
4. Extract text from each resulting file
5. Generate `extracted_content/extracted_content.json` with all results

## Output

The program creates:
- `extracted_content/` directory with split PDF files
- `extracted_content/extracted_content.json` containing:
  - Source PDF information
  - Individual file metadata
  - Extracted text content
  - Processing timestamps
  - File sizes and page counts

## JSON Structure

```json
[
  {
    "source_pdf": "Chords.pdf",
    "split_files": [
      {
        "file_name": "Chords_part_001.pdf",
        "extracted_at": "2025-06-26T21:06:32.925469-07:00",
        "text": "Extracted text content...",
        "page_count": 3,
        "file_size_bytes": 794773
      }
    ],
    "processing_time_seconds": 31026729958,
    "total_files_processed": 3
  }
]
```

## How It Works

1. **Text Detection**: Uses `pdftotext` to find pages containing "SPLITME"
2. **PDF Splitting**: Uses `pdftk` to split PDFs at detected pages
3. **Text Extraction**: 
   - First attempts direct text extraction with `pdftotext`
   - Falls back to OCR with `tesseract` if direct extraction fails
4. **JSON Generation**: Creates structured output with all metadata

## File Structure

```
music_sheets/
├── pdf_processor.go      # Main Go program
├── monitor_progress.sh   # Progress monitoring script
├── README.md            # This file
├── .gitignore           # Git ignore rules
└── v4/                  # Input PDF directory
    ├── Chords.pdf
    ├── Exercises.pdf
    ├── Licks.pdf
    ├── Scales.pdf
    ├── Songs A-F.pdf
    ├── Songs G-L.pdf
    ├── Songs M-S.pdf
    ├── Songs T-Z.pdf
    └── Theory.pdf
```

## Performance

- Processing time varies based on PDF size and complexity
- OCR processing is slower but handles image-based PDFs
- Large PDFs may take several minutes to process

## Troubleshooting

### Common Issues

1. **"command not found" errors**: Make sure all required tools are installed
2. **OCR quality issues**: Ensure PDFs are high resolution for better OCR results
3. **Memory issues**: Large PDFs may require more RAM

### Dependencies

- `pdftk-java`: PDF manipulation
- `poppler`: PDF text extraction and conversion
- `tesseract`: OCR text recognition
- `tesseract-lang`: Additional language support

## License

This project is open source. Feel free to modify and distribute.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

**Note**: The `extracted_content/` directory is excluded from version control due to its large size. Run the processor locally to generate the extracted content. 