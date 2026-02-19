package utils

import (
	"image"
	"io"
	"mime/multipart"
)

func IsImageByDecode(file multipart.File) bool {
	file.Seek(0, io.SeekStart)
	defer file.Seek(0, io.SeekStart)
	_, _, err := image.Decode(file)
	return err == nil
}
