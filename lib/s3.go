package lib

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AWSInfo ...
type AWSInfo struct {
	AwsRegion          string
	S3Bucket           string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

var awsSession *session.Session

// InitS3 ...
func InitS3(a AWSInfo) error {
	if a.AwsAccessKeyID == "" {
		return fmt.Errorf("AwsAccessKeyID is empty")
	}
	if a.AwsSecretAccessKey == "" {
		return fmt.Errorf("AwsSecretAccessKey is empty")
	}

	awsSession = session.New(&aws.Config{
		Region:      aws.String(a.AwsRegion),
		Credentials: credentials.NewStaticCredentials(a.AwsAccessKeyID, a.AwsSecretAccessKey, "")})
	if awsSession == nil {
		log.Fatal("Error occurred initializing aws session")
	}

	s3Session := s3.New(awsSession)

	listResults, err := s3Session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		awsSession = nil

		return fmt.Errorf("Failed to list buckets: %s", err)
	}

	foundBucket := false
	for _, b := range listResults.Buckets {
		log.Printf("Found bucket: %s", aws.StringValue(b.Name))
		if strings.EqualFold(aws.StringValue(b.Name), a.S3Bucket) {
			foundBucket = true
			//break
		}
	}

	if !foundBucket {
		createResult, err := s3Session.CreateBucket(
			&s3.CreateBucketInput{Bucket: &a.S3Bucket})
		if err != nil {
			awsSession = nil

			log.Printf("Failed to create bucket: %s", err)
			return err
		}
		if createResult == nil {
			log.Fatal("Create bucket result is empty")
		}
		log.Printf("Created bucket: %s", aws.StringValue(createResult.Location))

		if err := s3Session.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &a.S3Bucket}); err != nil {
			awsSession = nil

			log.Printf("Error waiting for bucket [%s] to exist: %s", a.S3Bucket, err)
			return err
		}
	}

	return nil
}

// PostS3 ...
func PostS3(bb bytes.Buffer, a AWSInfo, delaySuccess, delayError time.Duration) error {
	if awsSession == nil {
		return fmt.Errorf("AWS S3 Session is empty")
	}
	if bb.Len() == 0 {
		return fmt.Errorf("AWS S3 buffer is empty")
	}

	awsKey := time.Now().Format(time.RFC3339Nano) + ".gz"
	log.Printf("AWS S3 key: %s", awsKey)

	s3Uploader := s3manager.NewUploader(awsSession)

	for {
		r, err := s3Uploader.Upload(&s3manager.UploadInput{
			Body:   bytes.NewReader(bb.Bytes()),
			Bucket: aws.String(a.S3Bucket),
			Key:    aws.String(awsKey),
		})
		if err == nil {
			if r == nil {
				log.Fatal("Nil S3 result")
			}

			log.Printf("S3 Location: %s", r.Location)

			if delaySuccess > 0 {
				log.Printf("S3 delay on success %dms", delaySuccess)
				time.Sleep(delaySuccess)
			}

			break
		}

		log.Printf("Error posting to S3: %s", err)

		if delayError > 0 {
			log.Printf("S3 delay on error %dms", delayError)
			time.Sleep(delayError)
		}
	}

	return nil
}
