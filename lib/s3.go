package lib

import (
	"bytes"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AWSInfo ...
type AWSInfo struct {
	AwsRegion          string
	S3Bucket           string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

// PostS3 ...
func PostS3(bb bytes.Buffer, a AWSInfo) {
	if bb.Len() > 0 {
		uploader := s3manager.NewUploader(
			session.New(&aws.Config{
				Region:      aws.String(a.AwsRegion),
				Credentials: credentials.NewStaticCredentials(a.AwsAccessKeyID, a.AwsSecretAccessKey, "")}))

		awsKey := aws.String(time.Now().Format(time.RFC3339Nano) + ".gz")
		log.Printf("AWS S3 key: %s", *awsKey)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   bytes.NewReader(bb.Bytes()),
			Bucket: aws.String(a.S3Bucket),
			Key:    awsKey,
		})
		if err != nil {
			log.Printf("Error posting to S3 %s", err)
			return
		}

		if result == nil {
			log.Fatal("Nil S3 result")
		}
		log.Printf("S3 Location: %s", result.Location)
	}
}
