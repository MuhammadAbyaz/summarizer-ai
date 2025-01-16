package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lpernett/godotenv"
)

type S3Client struct {
	uploader   *s3manager.Uploader
	bucketName string
}

func NewS3Client() S3Client {
	godotenv.Load()
	s3Session := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(s3Session)

	return S3Client{
		uploader:   uploader,
		bucketName: os.Getenv("BUCKET_NAME"),
	}
}

func (u S3Client) UploadFile(key string, file []byte) (string, error) {
	contentType := http.DetectContentType(file)

	res, err := u.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file, %w", err)
	}
	return res.Location, nil
}
