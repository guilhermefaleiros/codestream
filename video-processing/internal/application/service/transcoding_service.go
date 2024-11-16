package service

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/video-processing/internal/application/gateway"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type TranscodingConfig struct {
	InputFolder            string
	OutputBaseFolder       string
	RemoteInputFolder      string
	RemoteOutputBaseFolder string
}

type TranscodingService struct {
	storageGateway gateway.StorageGateway
	config         TranscodingConfig
}

func (t *TranscodingService) Execute(ctx context.Context, videoID string) error {
	fileName := fmt.Sprintf("%s.mp4", videoID)

	inputFilePath, err := t.Download(ctx, fileName, t.config.InputFolder)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	err = t.Transcode(inputFilePath, videoID)
	if err != nil {
		return fmt.Errorf("failed to transcode file: %w", err)
	}

	err = t.Upload(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	err = t.Clear(videoID)
	if err != nil {
		return fmt.Errorf("failed to clear folders: %w", err)
	}

	fmt.Printf("Video transcoding completed for video ID: %s\n", videoID)
	return nil
}

func (t *TranscodingService) Clear(videoID string) error {
	outputFolderPath := fmt.Sprintf("%s/%s", t.config.OutputBaseFolder, videoID)
	inputFilePath := fmt.Sprintf("%s/%s.mp4", t.config.InputFolder, videoID)
	err := os.RemoveAll(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove input folder: %w", err)
	}

	err = os.RemoveAll(outputFolderPath)
	if err != nil {
		return fmt.Errorf("failed to remove output folder: %w", err)
	}

	fmt.Printf("Folders cleared for video ID: %s\n", videoID)
	return nil
}

func (t *TranscodingService) Transcode(filePath string, videoID string) error {
	outputFolder := fmt.Sprintf("%s/%s", t.config.OutputBaseFolder, videoID)
	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFilePath := filepath.Join(outputFolder, "output.mpd")

	cmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-map", "0",
		"-codec:v", "libx264",
		"-codec:a", "aac",
		"-f", "dash",
		"-seg_duration", "4",
		"-use_timeline", "1",
		"-use_template", "1",
		"-adaptation_sets", "id=0,streams=v id=1,streams=a",
		outputFilePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg error: %s\n%s", err, string(output))
	}

	fmt.Printf("MPEG-DASH files created at: %s\n", outputFolder)
	return nil
}

func (t *TranscodingService) Upload(ctx context.Context, videoID string) error {
	folder := fmt.Sprintf("%s/%s", t.config.OutputBaseFolder, videoID)
	remoteOutputFolder := fmt.Sprintf("%s/%s", t.config.RemoteOutputBaseFolder, videoID)
	files, err := os.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(files)) // Canal para capturar erros

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wg.Add(1)

		go func(file os.DirEntry) {
			defer wg.Done()

			filePath := filepath.Join(folder, file.Name())

			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", filePath, err)
				errChan <- fmt.Errorf("failed to read file %s: %w", filePath, err)
				return
			}

			contentType := mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			remotePath := filepath.Join(remoteOutputFolder, file.Name())

			err = t.storageGateway.Upload(ctx, fileBytes, remotePath, contentType)
			if err != nil {
				log.Printf("Failed to upload file %s to S3: %v", file.Name(), err)
				errChan <- fmt.Errorf("failed to upload file %s to S3: %w", file.Name(), err)
				return
			}

			log.Printf("File %s uploaded successfully to %s", file.Name(), remotePath)
		}(file)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		for e := range errChan {
			log.Printf("Error during upload: %v", e)
		}
		return fmt.Errorf("one or more uploads failed, check logs for details")
	}

	fmt.Printf("All files uploaded to: %s\n", remoteOutputFolder)
	return nil
}

func (t *TranscodingService) Download(ctx context.Context, fileName string, destinationFolder string) (string, error) {
	fileKey := fmt.Sprintf("%s/%s", t.config.RemoteInputFolder, fileName)
	file, err := t.storageGateway.Download(ctx, fileKey)
	if err != nil {
		fmt.Println("Error downloading file")
		panic(err)
	}
	err = os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory")
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	filePath := filepath.Join(destinationFolder, fileName)

	err = os.WriteFile(filePath, file, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("File downloaded to: %s\n", filePath)
	return filePath, nil

}

func NewTranscodingService(storageGateway gateway.StorageGateway, config TranscodingConfig) *TranscodingService {
	return &TranscodingService{storageGateway, config}
}
