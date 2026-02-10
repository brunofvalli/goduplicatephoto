- [x] Verify that the copilot-instructions.md file in the .github directory is created.

- [x] Clarify Project Requirements
	Go command-line application that:
	- Takes a directory path as a parameter
	- Traverses the directory tree recursively
	- Finds image files (jpg, png, bmp, gif, etc.)
	- Converts images to lower resolution for processing
	- Generates signatures for image comparison
	- Detects duplicate images (same content, different resolutions)
	- Moves lower-resolution duplicates to a separate directory for review

- [x] Scaffold the Project
	Created complete Go project structure with:
	- cmd/main.go: Command-line entry point
	- pkg/detector/detector.go: Main duplicate detection logic
	- pkg/signature/signature.go: Image signature generation
	- pkg/image/image.go: Image processing utilities
	- go.mod: Module dependencies
	- README.md: Comprehensive documentation

- [x] Customize the Project
	Implemented all core features:
	- Recursive directory traversal with filepath.WalkDir
	- Support for JPG, PNG, BMP, GIF, TIFF, WebP formats
	- Thumbnail-based image signature using MD5 hashing
	- Resolution-aware duplicate detection
	- Automatic movement of lower-resolution duplicates
	- Verbose logging and detailed statistics
	- Filename conflict handling in output directory

- [ ] Install Required Extensions

- [x] Compile the Project
	Successfully built: goduplicatephoto.exe
	All dependencies resolved:
	- github.com/disintegration/imaging v1.6.2
	- golang.org/x/image v0.14.0
	Build completed without errors

- [ ] Create and Run Task

- [ ] Launch the Project

- [ ] Ensure Documentation is Complete
