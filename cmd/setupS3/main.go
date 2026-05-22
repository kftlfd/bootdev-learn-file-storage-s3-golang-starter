package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/awsS3"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	s3Service, err := awsS3.NewS3Service()
	if err != nil {
		log.Fatal(err)
	}

	client := s3Service.Client
	bucket := s3Service.Bucket

	_, err = client.PutBucketPolicy(context.Background(), &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(getPolicy(bucket)),
	})
	if err != nil {
		log.Fatalf("Fail adding policy: %v", err)
	}

	fmt.Println("S3 bucket setup done:", bucket)
}

func getPolicy(bucket string) string {
	return fmt.Sprintf(`{
  "Version":"2012-10-17",
  "Statement":[
    {
      "Effect":"Allow",
      "Principal":"*",
      "Action":["s3:GetObject"],
      "Resource":["arn:aws:s3:::%s/*"]
    }
  ]
}`, bucket)
}
