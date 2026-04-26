package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
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
