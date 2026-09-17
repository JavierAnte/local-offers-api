package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	// maxImageDimension caps the long edge of a stored image. The app
	// already compresses/resizes before upload — this is a backstop so
	// storage size stays predictable regardless of what the client sends.
	maxImageDimension = 1600
	jpegQuality       = 80
)

type UploadService struct {
	uploadDir string
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

// Save decodes an uploaded image (auto-correcting EXIF orientation), scales
// it down to fit within maxImageDimension if needed, and re-encodes it as a
// JPEG on disk under uploadDir. It returns the stored filename.
func (s *UploadService) Save(file io.Reader) (string, error) {
	img, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		return "", ErrInvalidInput
	}

	resized := imaging.Fit(img, maxImageDimension, maxImageDimension, imaging.Lanczos)

	filename := uuid.New().String() + ".jpg"

	err = os.MkdirAll(s.uploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	out, err := os.Create(filepath.Join(s.uploadDir, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	err = imaging.Encode(out, resized, imaging.JPEG, imaging.JPEGQuality(jpegQuality))
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}
