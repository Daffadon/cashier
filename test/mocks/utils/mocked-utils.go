package test

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockFileManagement struct {
	mock.Mock
}

func (m *MockFileManagement) UploadFile(file *multipart.FileHeader, filename, path string) error {
	args := m.Called(file, filename, path)
	return args.Error(0)
}
func (m *MockFileManagement) GetFileNameExtension(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}
func (m *MockFileManagement) DeleteFile(pathFile string) error {
	args := m.Called(pathFile)
	return args.Error(0)
}

func (m *MockFileManagement) GenerateNewFileName(ext string) string {
	args := m.Called(ext)
	return args.String(0)
}
