package signature

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/brunofvalli/goduplicatephoto/pkg/image"
)

// GenerateSignature creates a hash-based signature for an image
// It resizes the image to a thumbnail and hashes the pixel data
// This allows detection of duplicate images at different resolutions
func GenerateSignature(filePath string, thumbSize int) (string, error) {
	// Open the image file
	img, err := image.LoadImage(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to load image: %w", err)
	}

	// Create a thumbnail (normalized size for comparison)
	thumb := image.CreateThumbnail(img, thumbSize)

	// Convert thumbnail to byte data
	thumbData := image.ImageToBytes(thumb)

	// Generate MD5 hash of the thumbnail data
	hash := md5.New()
	_, err = hash.Write(thumbData)
	if err != nil {
		return "", fmt.Errorf("failed to compute hash: %w", err)
	}

	signature := fmt.Sprintf("%x", hash.Sum(nil))
	return signature, nil
}

// GenerateFileHash creates a hash of the entire file
// Useful for quick detection of identical files
func GenerateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute file hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
