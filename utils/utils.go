package utils

import (
	"encoding/base64"
	_ "golang.org/x/image/webp" // Register WEBP format
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io/fs"
	"os"
	"strings"
)

func CreateFolder(folderPath string, perm fs.FileMode) error {
	var err error
	if _, err = os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

func Base64ToImage(dataBase64 string) (img image.Image, format string, err error) {
	// Create a reader from the decoded byte array
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(dataBase64))

	// Decode the image to get the dimensions
	return image.Decode(reader)
}
func SaveFileLocal(filepath string, data []byte) error {
	var err error
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return err
}

func DeleteFileLocal(filepath string) error {
	var err error
	err = os.Remove(filepath)
	if err != nil {
		return err
	}

	return err
}
