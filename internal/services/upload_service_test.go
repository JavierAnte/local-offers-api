package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadServiceSave(t *testing.T) {
	t.Parallel()

	var input bytes.Buffer
	source := image.NewRGBA(image.Rect(0, 0, 2000, 1000))
	source.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&input, source); err != nil {
		t.Fatalf("encode fixture: %v", err)
	}

	uploadDir := t.TempDir()
	filename, err := NewUploadService(uploadDir).Save(&input)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if !strings.HasSuffix(filename, ".jpg") {
		t.Fatalf("filename = %q, want .jpg suffix", filename)
	}

	storedPath := filepath.Join(uploadDir, filename)
	stored, err := os.Open(storedPath)
	if err != nil {
		t.Fatalf("open stored image: %v", err)
	}
	defer stored.Close()

	decoded, err := jpeg.Decode(stored)
	if err != nil {
		t.Fatalf("stored file is not JPEG: %v", err)
	}
	if got := decoded.Bounds().Size(); got.X != 1600 || got.Y != 800 {
		t.Fatalf("stored dimensions = %dx%d, want 1600x800", got.X, got.Y)
	}
}

func TestUploadServiceRejectsInvalidImage(t *testing.T) {
	t.Parallel()

	_, err := NewUploadService(t.TempDir()).Save(strings.NewReader("not an image"))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Save() error = %v, want ErrInvalidInput", err)
	}
}

func TestUploadServiceReportsDirectoryFailure(t *testing.T) {
	t.Parallel()

	var input bytes.Buffer
	if err := png.Encode(&input, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatalf("encode fixture: %v", err)
	}

	notDirectory := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(notDirectory, []byte("occupied"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := NewUploadService(notDirectory).Save(&input); err == nil {
		t.Fatal("Save() error = nil, want directory error")
	}
}
