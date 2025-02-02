package fileManager

import (
	"github.com/sirupsen/logrus"
	"os"
	"styl-monolith/internal/media/core/ports"
)

type FileManager struct {
	log *logrus.Logger
}

func NewFileManager(log *logrus.Logger) ports.FileManager {
	return &FileManager{
		log: log,
	}
}

// CheckIfFolderCreated checks if a folder exists at the given path.
func (fm *FileManager) CheckIfFolderCreated(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// CreateFolder creates a folder at the specified path.
func (fm *FileManager) CreateFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFolder deletes a folder at the specified path.
func (fm *FileManager) DeleteFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

// DeleteImage deletes a file (image) at the specified path.
func (fm *FileManager) DeleteImage(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// GetImage retrieves the content of a file (image) as a byte slice.
func (fm *FileManager) GetImage(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveImage saves a byte slice as a file (image) at the specified path.
func (fm *FileManager) SaveImage(path string, data []byte) error {
	err := os.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
