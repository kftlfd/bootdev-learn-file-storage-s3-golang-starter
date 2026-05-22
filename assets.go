package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
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
