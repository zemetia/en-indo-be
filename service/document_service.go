package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DocumentService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
	isAllowedExt(ext string) bool
}

type documentService struct {
	baseURL     string
	uploadPath  string
	allowedExts []string
}

type DocumentConfig struct {
	// ... existing fields ...
	BaseURL    string `mapstructure:"BASE_URL"`
	UploadPath string `mapstructure:"UPLOAD_PATH"`
}

func NewDocumentService(baseURL, uploadPath string) DocumentService {
	return &documentService{
		baseURL:    baseURL,
		uploadPath: uploadPath,
		allowedExts: []string{
			".jpg", ".jpeg", ".png", ".gif",
		},
	}
}

func (s *documentService) UploadImage(file *multipart.FileHeader) (string, error) {
	// Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedExt(ext) {
		return "", errors.New("file extension not allowed")
	}

	// Validasi ukuran file (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", errors.New("file size too large (max 5MB)")
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Buat folder jika belum ada
	uploadDir := filepath.Join(s.uploadPath, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Path lengkap file
	filepath := filepath.Join(uploadDir, filename)

	// Buka file yang diupload
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Buat file baru di server
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file
	buffer := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := src.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", fmt.Errorf("failed to read uploaded file: %v", err)
		}

		if _, err := dst.Write(buffer[:n]); err != nil {
			return "", fmt.Errorf("failed to write file: %v", err)
		}
	}

	// Return URL file
	relativePath := strings.TrimPrefix(filepath, s.uploadPath)
	return fmt.Sprintf("%s/uploads%s", s.baseURL, relativePath), nil
}

func (s *documentService) isAllowedExt(ext string) bool {
	for _, allowedExt := range s.allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
