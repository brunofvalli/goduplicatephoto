# Go Duplicate Photo Detector

A command-line application written in Go that finds and manages duplicate photos across a directory tree.

## Features

- **Directory Traversal**: Recursively scans a directory and all subdirectories
- **Multi-format Support**: Supports JPG, PNG, BMP, GIF, TIFF, and WebP formats
- **Resolution-aware**: Creates image signatures that detect duplicates regardless of resolution
- **Smart Signature Generation**: Uses thumbnail-based hashing to compare images efficiently
- **Duplicate Detection**: Identifies duplicate images even if they have different resolutions
- **Automatic Organization**: Moves lower-resolution duplicates to a designated directory for review
- **Detailed Reporting**: Provides statistics on scanned files, images found, and duplicates detected

## Installation

### Prerequisites

- Go 1.21 or later
- Image processing libraries (installed via `go mod download`)

### Build

```bash
go build -o goduplicatephoto ./cmd
```

### Dependencies

The project uses:
- `github.com/disintegration/imaging` - Image processing
- `github.com/dsoprea/go-exif/v3` - EXIF data handling

Run `go mod download` to install dependencies.

## Usage

### Basic Usage

```bash
goduplicatephoto -dir /path/to/images
```

### Command-line Options

- `-dir` (string): Directory to scan for duplicate photos (default: ".")
- `-output` (string): Output directory for duplicate photos (default: "duplicates" subfolder in input dir)
- `-thumbsize` (int): Thumbnail size for signature generation in pixels (default: 200)
- `-verbose`: Enable verbose output for detailed processing information

### Examples

```bash
# Scan current directory
goduplicatephoto

# Scan specific directory
goduplicatephoto -dir C:\Pictures

# Scan with custom output directory
goduplicatephoto -dir C:\Pictures -output C:\Pictures\Duplicates_Review

# Scan with verbose output and custom thumbnail size
goduplicatephoto -dir /home/user/photos -verbose -thumbsize 256
```

## How It Works

1. **Scanning**: Walks through the entire directory tree, identifying image files
2. **Signature Generation**: 
   - Loads each image
   - Creates a thumbnail (normalized to a standard size)
   - Generates an MD5 hash of the thumbnail data
3. **Duplicate Detection**: Groups images with identical signatures
4. **Resolution Sorting**: Within each duplicate group, sorts images by resolution
5. **File Organization**: Moves lower-resolution duplicates to the output directory

## Output

The application provides a summary report:

```
--- Duplicate Photo Detection Summary ---
Total files scanned: 1523
Valid images found: 342
Duplicates found: 28
Files moved: 47
Output directory: C:\Pictures\duplicates
```

## Project Structure

```
scr/
├── cmd/
│   └── main.go              # Entry point
├── pkg/
│   ├── detector/
│   │   └── detector.go      # Main detection logic
│   ├── signature/
│   │   └── signature.go     # Image signature generation
│   └── image/
│       └── image.go         # Image processing utilities
├── go.mod
├── go.sum
└── README.md
```

## Important Notes

- **Backup First**: Always backup your photos before running this tool
- **Review Duplicates**: The application moves suspected duplicates to a separate folder - review them before permanent deletion
- **Permissions**: Ensure you have read/write permissions in the scanned directories
- **Performance**: Large image collections may take time to process
- **Duplicate Definition**: Images are considered duplicates if they have the same thumbnail signature, regardless of resolution or metadata

## Future Enhancements

- Perceptual hashing for more robust duplicate detection
- Parallel processing for faster scanning
- Web interface for preview and management
- EXIF-based metadata comparison
- Dry-run mode to preview changes without moving files
- Database of processed files to avoid re-scanning

## License

MIT License

## Author

Bruno Valli
