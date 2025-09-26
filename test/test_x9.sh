#!/bin/bash

# X9 File Testing Script with Null Padding Detection
# Usage: ./test_x9.sh <input_file>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <input_file>"
    echo "Example: $0 /path/to/your/file.x937"
    exit 1
fi

INPUT_FILE="$1"
CLEANED_FILE="/tmp/x9_cleaned_$(date +%s).x937"

echo "🔍 X9 File Testing with Null Padding Detection"
echo "=============================================="
echo "Input file: $INPUT_FILE"

# Check if file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "❌ Error: File does not exist: $INPUT_FILE"
    exit 1
fi

# Step 1: Handle null padding
echo ""
echo "Step 1: Checking for null padding..."
go run simple_strip.go "$INPUT_FILE" "$CLEANED_FILE"

if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to process file"
    exit 1
fi

	# Step 2: Test the file
	echo ""
	echo "Step 2: Testing X9 file..."
	echo "=========================="

	# Set FRB compatibility mode for permissive validation
	export FRB_COMPATIBILITY_MODE="true"
	
	go run validate_x9.go "$CLEANED_FILE"

# Check the result
if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 SUCCESS: File parsed completely!"
    echo "   • All validation checks passed"
    echo "   • File structure is valid"
else
    echo ""
    echo "ℹ️  File had validation errors:"
    echo "   • Check the error messages above for specific field issues"
    echo "   • Common issues: missing required fields, invalid field values"
fi

# Cleanup
rm -f "$CLEANED_FILE"
echo ""
echo "🧹 Cleaned up temporary file"
echo ""
echo "📝 Summary:"
echo "   • Null padding detection: ✅ Working"
echo "   • File processing: ✅ Complete"
echo "   • Validation: See results above"
