#!/bin/bash

# PDF Metadata Modifier Script
# This script demonstrates how to modify PDF metadata using exiftool
# Note: Requires exiftool to be installed

# Check if exiftool is installed
if ! command -v exiftool &> /dev/null; then
    echo "Error: exiftool is not installed."
    echo "Install it with: brew install exiftool (macOS) or apt-get install exiftool (Ubuntu)"
    exit 1
fi

# Function to process a single file
process_file() {
    local source_file="$1"
    local output_file="$2"
    local title="$3"
    local author="$4"
    local subject="$5"
    local keywords="$6"
    
    echo "Processing: $source_file -> $output_file"
    
    # Copy the file first
    cp "$source_file" "$output_file"
    
    # Modify metadata using exiftool
    exiftool -overwrite_original \
        -Title="$title" \
        -Author="$author" \
        -Subject="$subject" \
        -Keywords="$keywords" \
        "$output_file" > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "  ✓ Metadata updated successfully"
    else
        echo "  ✗ Failed to update metadata"
    fi
}

# Example usage
echo "PDF Metadata Modifier Example"
echo "============================="

# Create output directory
mkdir -p processed_files_with_metadata

# Example: Process a single file
if [ -f "processed_files/Blue_Monk_Jazz.pdf" ]; then
    process_file \
        "processed_files/Blue_Monk_Jazz.pdf" \
        "processed_files_with_metadata/Blue_Monk_Jazz_with_metadata.pdf" \
        "Blue Monk" \
        "Thelonious Monk" \
        "Jazz" \
        "song, lead-sheet"
else
    echo "Example file not found. Run the Go program first to generate processed files."
fi

echo ""
echo "To process all files, you would need to:"
echo "1. Parse the metadata_clean.json file"
echo "2. Call this script for each file with the appropriate metadata"
echo "3. Or use a tool like pdftk for batch processing"
echo ""
echo "Alternative command-line tools:"
echo "- pdftk: pdftk input.pdf update_info metadata.txt output output.pdf"
echo "- qpdf: qpdf --replace-input --update-info-from=metadata.txt input.pdf" 