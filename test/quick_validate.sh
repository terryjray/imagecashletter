#!/bin/bash

# Quick X9 File Validation Script
# Usage: ./quick_validate.sh <input_file>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <input_file>"
    echo "Example: $0 /path/to/your/file.x937"
    exit 1
fi

INPUT_FILE="$1"

echo "🔍 Quick X9 File Validation"
echo "=========================="
echo "Input file: $INPUT_FILE"

# Check if file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "❌ Error: File does not exist: $INPUT_FILE"
    exit 1
fi

# Set FRB compatibility mode for permissive validation
export FRB_COMPATIBILITY_MODE="true"

# Run validation
echo ""
go run validate_x9.go "$INPUT_FILE"

# Check the result
if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 SUCCESS: File is valid and ready for processing!"
else
    echo ""
    echo "ℹ️  File has issues - check the error messages above"
    echo "   • Try running with preprocessing: ./test_x9.sh $INPUT_FILE"
fi
