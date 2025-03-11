package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"tiga-putra-cashier-be/dto"
)

func GetFileNameExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return strings.ToLower(parts[len(parts)-1])
}

func UploadFile(file *multipart.FileHeader, filename, path string) error {

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

func DeleteFile(pathFile string) error {
	if err := os.Remove(pathFile); err != nil {
		return err
	}
	return nil
}
