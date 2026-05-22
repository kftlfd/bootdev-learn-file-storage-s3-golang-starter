package awsS3

import (
	"fmt"
	"os"
)

type s3Config struct {
	Url            string
	Bucket         string
	Region         string
	CfDistribution string
}

func loadS3ConfigFromEnv() (s3Config, error) {
	cfg := s3Config{}

	cfg.Url = os.Getenv("S3_URL")
	if cfg.Url == "" {
		return cfg, fmt.Errorf("S3_URL environment variable is not set")
	}

	cfg.Bucket = os.Getenv("S3_BUCKET")
	if cfg.Bucket == "" {
		return cfg, fmt.Errorf("S3_BUCKET environment variable is not set")
	}

	cfg.Region = os.Getenv("S3_REGION")
	if cfg.Region == "" {
		return cfg, fmt.Errorf("S3_REGION environment variable is not set")
	}

	cfg.CfDistribution = os.Getenv("S3_CF_DISTRO")
	if cfg.CfDistribution == "" {
		return cfg, fmt.Errorf("S3_CF_DISTRO environment variable is not set")
	}

	return cfg, nil
}
