package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/brunofvalli/goduplicatephoto/pkg/image"
	"github.com/brunofvalli/goduplicatephoto/pkg/signature"
)

// Config holds the configuration for duplicate detection
type Config struct {
	InputDir      string
	OutputDir     string
	ThumbnailSize int
	Verbose       bool
}

// Stats holds statistics about the detection process
type Stats struct {
	TotalFiles      int
	ImagesFound     int
	DuplicatesFound int
	FilesMoved      int
}

// DuplicateDetector handles the detection of duplicate photos
type DuplicateDetector struct {
	config    *Config
	sigs      map[string][]string // signature -> list of file paths
	mu        sync.Mutex
}

// NewDuplicateDetector creates a new duplicate detector
func NewDuplicateDetector(config *Config) *DuplicateDetector {
	return &DuplicateDetector{
		config: config,
		sigs:   make(map[string][]string),
	}
}

// Run performs the duplicate detection process
func (dd *DuplicateDetector) Run() (*Stats, error) {
	stats := &Stats{}

	// Walk the directory tree
	err := filepath.WalkDir(dd.config.InputDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		stats.TotalFiles++

		// Check if file is an image
		if !image.IsImageFile(path) {
			return nil
		}

		stats.ImagesFound++

		// Generate signature for the image
		sig, err := signature.GenerateSignature(path, dd.config.ThumbnailSize)
		if err != nil {
			if dd.config.Verbose {
				fmt.Printf("Warning: Failed to process image %s: %v\n", path, err)
			}
			return nil
		}

		// Store signature and file path
		dd.mu.Lock()
		dd.sigs[sig] = append(dd.sigs[sig], path)
		dd.mu.Unlock()

		if dd.config.Verbose {
			fmt.Printf("Processed: %s (sig: %s)\n", path, sig[:8]+"...")
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory tree: %w", err)
	}

	// Find and move duplicates
	moved, err := dd.processDuplicates()
	if err != nil {
		return nil, err
	}

	stats.FilesMoved = moved
	stats.DuplicatesFound = dd.countDuplicates()

	return stats, nil
}

// processDuplicates handles the movement of duplicate files
func (dd *DuplicateDetector) processDuplicates() (int, error) {
	moved := 0

	for sig, files := range dd.sigs {
		if len(files) <= 1 {
			continue // Not a duplicate
		}

		if dd.config.Verbose {
			fmt.Printf("Found %d duplicates with signature %s\n", len(files), sig[:8]+"...")
		}

		// Sort files by resolution/file size to keep highest quality
		sortedFiles, err := image.SortByResolution(files)
		if err != nil {
			return moved, fmt.Errorf("error sorting files: %w", err)
		}

		// Move lower resolution duplicates to output directory
		for i := 1; i < len(sortedFiles); i++ {
			destPath := filepath.Join(dd.config.OutputDir, filepath.Base(sortedFiles[i]))

			// Create output directory if needed
			if err := os.MkdirAll(dd.config.OutputDir, 0755); err != nil {
				return moved, fmt.Errorf("error creating output directory: %w", err)
			}

			// Handle filename conflicts
			destPath = handleDuplicateFilename(destPath)

			if err := os.Rename(sortedFiles[i], destPath); err != nil {
				return moved, fmt.Errorf("error moving file %s: %w", sortedFiles[i], err)
			}

			if dd.config.Verbose {
				fmt.Printf("Moved: %s -> %s\n", sortedFiles[i], destPath)
			}

			moved++
		}
	}

	return moved, nil
}

// countDuplicates returns the number of duplicate signatures found
func (dd *DuplicateDetector) countDuplicates() int {
	count := 0
	for _, files := range dd.sigs {
		if len(files) > 1 {
			count++
		}
	}
	return count
}

// handleDuplicateFilename handles filename conflicts in the output directory
func handleDuplicateFilename(path string) string {
	if _, err := os.Stat(path); err == nil {
		// File exists, add a number
		dir := filepath.Dir(path)
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		nameWithoutExt := base[:len(base)-len(ext)]

		for i := 1; i < 1000; i++ {
			newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", nameWithoutExt, i, ext))
			if _, err := os.Stat(newPath); err != nil {
				return newPath
			}
		}
	}
	return path
}
