package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var videoMeta ffprobeOutput
	err := json.Unmarshal(stdout.Bytes(), &videoMeta)
	if err != nil {
		return "", err
	}

	if len(videoMeta.Streams) == 0 {
		return "", fmt.Errorf("no streams found in videoMetadata")
	}

	width := videoMeta.Streams[0].Width
	height := videoMeta.Streams[0].Height
	ratio := float64(width) / float64(height)
	const tolerance = 0.01
	switch {
	case math.Abs(ratio-(16.0/9.0)) < tolerance:
		return "16:9", nil
	case math.Abs(ratio-(9.0/16.0)) < tolerance:
		return "9:16", nil
	default:
		return "other", nil
	}

}

func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilePath)
	// Capturing error details (stderr)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg failed: %v, stderr: %s", err, stderr.String())
	}

	// Stat-checking the output
	fileInfo, err := os.Stat(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("could not stat processed file: %v", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("processed file is empty")
	}

	return outputFilePath, nil
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignedUrl, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}
	return presignedUrl.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(dbVideo database.Video) (database.Video, error) {
	if dbVideo.VideoURL == nil {
		return dbVideo, nil
	}

	parts := strings.SplitN(*dbVideo.VideoURL, ",", 2)
	if len(parts) != 2 {
		return dbVideo, nil
	}

	bucket := parts[0]
	key := parts[1]
	presigned, err := generatePresignedURL(cfg.s3Client, bucket, key, 5*time.Minute)
	if err != nil {
		return dbVideo, err
	}

	dbVideo.VideoURL = &presigned
	return dbVideo, nil
}
