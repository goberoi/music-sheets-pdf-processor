#!/bin/bash

echo "Monitoring PDF processing progress..."
echo "====================================="

while true; do
    echo "$(date): $(ls extracted_content/*.pdf 2>/dev/null | wc -l | tr -d ' ') PDF files created"
    
    if [ -f "extracted_content/extracted_content.json" ]; then
        echo "✅ Processing complete! JSON file created."
        echo "📊 Summary:"
        echo "   Total files: $(ls extracted_content/*.pdf 2>/dev/null | wc -l | tr -d ' ')"
        echo "   JSON file size: $(ls -lh extracted_content/extracted_content.json | awk '{print $5}')"
        break
    fi
    
    sleep 10
done 