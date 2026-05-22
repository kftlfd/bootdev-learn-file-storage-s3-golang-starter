package awsS3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newS3Client(cfg s3Config) (*s3.Client, error) {
	s3ClientConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(s3ClientConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Url)
		o.UsePathStyle = true
	})

	return s3Client, nil
}
