package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createRandomAssetKey(mediaType string, suffix ...string) string {
	var fileName strings.Builder

	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	fileName.WriteString(base64.RawURLEncoding.EncodeToString(base))

	for _, suf := range suffix {
		fileName.WriteString(".")
		fileName.WriteString(suf)
	}

	fileName.WriteString(".")
	fileName.WriteString(mediaTypeToExt(mediaType))

	return fileName.String()
}

func mediaTypeToExt(mt string) string {
	_, ext, ok := strings.Cut(mt, "/")
	if !ok || len(ext) < 1 {
		return "bin"
	}
	return ext
}

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getDiskObkectUrl(key string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, key)
}

func (cfg apiConfig) getS3ObjectUrl(key string) string {
	return fmt.Sprintf("%s/%s/%s", cfg.s3.Url, cfg.s3.Bucket, key)
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error running ffprobe: %w", err)
	}

	res := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{}
	if err := json.NewDecoder(buf).Decode(&res); err != nil {
		return "", fmt.Errorf("Error decoding ffprobe result: %w", err)
	}

	ratio := "other"

	const tolerance = 0.01
	r := float64(res.Width) / float64(res.Height)
	if math.Abs(r-(16.0/9.0)) <= tolerance {
		ratio = "16:9"
	} else if math.Abs(r-(9.0/16.0)) <= tolerance {
		ratio = "9:16"
	}

	return ratio, nil
}

func getS3ObjVideoPrefix(filePath string) (string, error) {
	ratio, err := getVideoAspectRatio(filePath)
	if err != nil {
		return "", err
	}

	if ratio == "16:9" {
		return "landscape", nil
	}
	if ratio == "9:16" {
		return "portrait", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outPath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outPath)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error running ffmpeg: %w", err)
	}

	return outPath, nil
}
