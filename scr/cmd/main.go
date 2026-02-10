package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/brunofvalli/goduplicatephoto/pkg/detector"
)

func main() {
	var (
		inputDir  = flag.String("dir", ".", "Directory to scan for duplicate photos")
		outputDir = flag.String("output", "", "Output directory for duplicate photos (if empty, creates 'duplicates' in input dir)")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
		thumbSize = flag.Int("thumbsize", 200, "Thumbnail size for signature generation (pixels)")
	)

	flag.Parse()

	if *inputDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate input directory
	info, err := os.Stat(*inputDir)
	if err != nil {
		log.Fatalf("Error accessing directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Input path is not a directory: %s", *inputDir)
	}

	// Set output directory
	if *outputDir == "" {
		*outputDir = filepath.Join(*inputDir, "duplicates")
	}

	if *verbose {
		fmt.Printf("Input Directory: %s\n", *inputDir)
		fmt.Printf("Output Directory: %s\n", *outputDir)
		fmt.Printf("Thumbnail Size: %d\n", *thumbSize)
		fmt.Println("Starting duplicate photo detection...")
	}

	// Create detector and run analysis
	config := &detector.Config{
		InputDir:      *inputDir,
		OutputDir:     *outputDir,
		ThumbnailSize: *thumbSize,
		Verbose:       *verbose,
	}

	detector := detector.NewDuplicateDetector(config)
	stats, err := detector.Run()
	if err != nil {
		log.Fatalf("Error during detection: %v", err)
	}

	// Print results
	fmt.Println("\n--- Duplicate Photo Detection Summary ---")
	fmt.Printf("Total files scanned: %d\n", stats.TotalFiles)
	fmt.Printf("Valid images found: %d\n", stats.ImagesFound)
	fmt.Printf("Duplicates found: %d\n", stats.DuplicatesFound)
	fmt.Printf("Files moved: %d\n", stats.FilesMoved)
	if stats.DuplicatesFound > 0 {
		fmt.Printf("Output directory: %s\n", *outputDir)
	}
}
