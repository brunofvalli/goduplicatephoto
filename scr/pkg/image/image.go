package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"github.com/disintegration/imaging"
)

// Supported image extensions
var supportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".bmp":  true,
	".gif":  true,
	".tiff": true,
	".webp": true,
}

// IsImageFile checks if a file is a supported image format
func IsImageFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	return supportedExtensions[ext]
}

// LoadImage loads an image from a file path
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// CreateThumbnail creates a thumbnail of the image
func CreateThumbnail(img image.Image, size int) image.Image {
	return imaging.Thumbnail(img, size, size, imaging.Lanczos)
}

// ImageToBytes converts an image to a byte array for hashing
func ImageToBytes(img image.Image) []byte {
	var buf bytes.Buffer

	// Encode image as PNG to get consistent byte representation
	err := png.Encode(&buf, img)
	if err != nil {
		// Fallback to JPEG if PNG encoding fails
		_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	}

	return buf.Bytes()
}

// GetImageDimensions returns the width and height of an image file
func GetImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	return img.Width, img.Height, nil
}

// SortByResolution sorts image files by their resolution (descending)
// Files with higher resolution come first
func SortByResolution(filePaths []string) ([]string, error) {
	type fileInfo struct {
		path     string
		width    int
		height   int
		megapixels float64
	}

	var files []fileInfo

	for _, path := range filePaths {
		width, height, err := GetImageDimensions(path)
		if err != nil {
			continue
		}

		megapixels := float64(width*height) / 1_000_000
		files = append(files, fileInfo{
			path:       path,
			width:      width,
			height:     height,
			megapixels: megapixels,
		})
	}

	// Sort by megapixels (descending)
	sort.Slice(files, func(i, j int) bool {
		return files[i].megapixels > files[j].megapixels
	})

	// Extract sorted paths
	sorted := make([]string, len(files))
	for i, f := range files {
		sorted[i] = f.path
	}

	return sorted, nil
}

// ConvertToGrayscale converts an image to grayscale for signature comparison
func ConvertToGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert RGB to grayscale using standard formula
			gray.Set(x, y, color.Gray{Y: uint8((r + g + b) / 3)})
		}
	}

	return gray
}
