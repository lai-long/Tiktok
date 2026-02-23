package utils

import (
	"mime/multipart"
	"net/http"
)

func IsImage(file multipart.File) (bool, error) {
	head := make([]byte, 512)
	_, err := file.Read(head)
	if err != nil {
		return false, err
	}
	mime := http.DetectContentType(head)
	switch mime {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true, nil
	default:
		return false, nil
	}
}
