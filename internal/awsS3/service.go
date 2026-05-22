package awsS3

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3Service struct {
	s3Config
	Client *s3.Client
}

func NewS3Service() (S3Service, error) {
	cfg, err := loadS3ConfigFromEnv()
	if err != nil {
		return S3Service{}, err
	}

	client, err := newS3Client(cfg)
	if err != nil {
		return S3Service{}, err
	}

	return S3Service{
		s3Config: cfg,
		Client:   client,
	}, nil
}
