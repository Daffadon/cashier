package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"tiga-putra-cashier-be/dto"

	"github.com/google/uuid"
)

type FileManagement interface {
	UploadFile(file *multipart.FileHeader, filename, path string) error
	GetFileNameExtension(filename string) string
	DeleteFile(pathFile string) error
	GenerateNewFileName(ext string) string
}
type fileManagementUtils struct{}

func FileInit() FileManagement {
	return &fileManagementUtils{}
}

func (f *fileManagementUtils) GetFileNameExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return strings.ToLower(parts[len(parts)-1])
}

func (f *fileManagementUtils) UploadFile(file *multipart.FileHeader, filename, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0750); err != nil {
			return dto.ErrToSaveFile
		}
	}
	fullFilePath := fmt.Sprintf("%s/%s", path, filename)
	filepath.Clean(fullFilePath)
	uploadedFile, err := file.Open()
	if err != nil {
		return dto.ErrToSaveFile
	}
	defer uploadedFile.Close()
	targetFile, err := os.Create(fullFilePath)
	if err != nil {
		return dto.ErrToSaveFile
	}

	defer targetFile.Close()
	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		return dto.ErrToSaveFile
	}
	return nil
}

func (f *fileManagementUtils) DeleteFile(pathFile string) error {
	if err := os.Remove(pathFile); err != nil {
		return err
	}
	return nil
}

func (f *fileManagementUtils) GenerateNewFileName(ext string) string {
	return fmt.Sprintf("%s.%s", uuid.New().String(), ext)
}
