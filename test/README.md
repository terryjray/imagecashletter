# X9 File Test Utilities

This directory contains utilities for validating X9 image cash letter files, designed for support technicians to quickly diagnose file issues.

## Files

- `validate_x9.go` - **Main validation utility** - automatically detects encoding and validates files
- `quick_validate.sh` - **Simple validation script** - quick validation without preprocessing
- `test_x9.sh` - **Complete testing script** - includes null padding detection and preprocessing
- `simple_strip.go` - Standalone utility to detect and strip null padding (for reference)

## Quick Start

### Step 1: Try Quick Validation (Most Files)
```bash
./quick_validate.sh /path/to/your/file.x937
```
**Use this first** - works for 80% of files (clean, well-formatted files)

### Step 2: If Quick Validation Fails, Try Preprocessing
```bash
./test_x9.sh /path/to/your/file.x937
```
**Use this when** quick validation fails - handles files with:
- Null padding prefixes
- Encoding issues  
- Structural problems
- Vendor format variations

## Features

### Automatic Encoding Detection
- **ASCII**: Detects standard ASCII-encoded X9 files
- **EBCDIC**: Automatically detects and handles EBCDIC-encoded files
- **Reader Options**: Applies appropriate reader options based on encoding

### Error Diagnosis
- **Specific Error Messages**: Provides targeted guidance for common issues
- **Vendor Variations**: Handles different file formats from various vendors
- **Null Padding**: Automatically detects and handles null padding prefixes

### File Summary
- **Cash Letter Count**: Number of cash letters in the file
- **Record Count**: Total number of records
- **Item Count**: Number of check items
- **Total Amount**: Dollar amount of all items
- **Business Details**: Cash letter ID, dates, routing numbers

## Example Output

```
🔍 Validating X9 file: /path/to/file.x937
============================================================
📝 Detected encoding: EBCDIC
🔧 EBCDIC encoding option enabled

📖 Reading file...
✅ File read successfully!

🔍 Validating file structure...
✅ File validation passed!

📊 File Summary:
   • Cash Letters: 1
   • Total Records: 319
   • Total Items: 31
   • Total Amount: $824569.94

📋 First Cash Letter Details:
   • ID: '00000206'
   • Business Date: 2025-09-26
   • Destination: 062203308
   • Originator: Operational Se

🎉 File is valid and ready for processing!
```

## Troubleshooting

### Common Error Messages

| Error | Solution |
|-------|----------|
| `unknown record type` | Try `./test_x9.sh` (file has null padding) |
| `TestFileIndicator is mandatory` | Try `./test_x9.sh` (file needs preprocessing) |
| `unrecognized encoding` | File format issue - check with vendor |
| `token too long` | File structure issue - try preprocessing |

### Workflow
1. **Always start with**: `./quick_validate.sh file.x937`
2. **If it fails**: `./test_x9.sh file.x937` 
3. **If both fail**: Check file format with vendor

## Notes

- **Encoding detection is automatic** - no need to specify ASCII vs EBCDIC
- **Error messages provide specific guidance** for common vendor issues
- **File summaries help verify** the file contains expected data
- **Both preprocessing and direct validation** options available
