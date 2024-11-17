package pkg

import (
	"fmt"
	"net/http"
)

func GetFileContentType(file []byte) string {
	return http.DetectContentType(file)
}

func GetMimeToExtension(contentType string) (string, error) {
	mimeToExt := map[string]string{
		"image/jpeg":       ".jpg",
		"image/png":        ".png",
		"image/gif":        ".gif",
		"video/mp4":        ".mp4",
		"video/x-matroska": ".mkv",
		"audio/mpeg":       ".mp3",
		"application/pdf":  ".pdf",
		"text/plain":       ".txt",
	}
	ext, ok := mimeToExt[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
	return ext, nil
}
