package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func JSON(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func JSONError(w http.ResponseWriter, message string, code int) {
	JSON(w, map[string]string{"error": message}, code)
}

func JSONWithLocation(w http.ResponseWriter, response interface{}, location string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().
		Set("Location", location)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func DeserializeRequestBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func FileToBytes(file multipart.File) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return buf.Bytes(), nil
}

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
