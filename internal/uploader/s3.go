package uploader

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Upload encodes the image as PNG and uploads it to S3.
func Upload(ctx context.Context, bucket, key string, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(buf.Bytes()),
		ContentType:  aws.String("image/png"),
		CacheControl: aws.String("max-age=300"),
	})
	if err != nil {
		return fmt.Errorf("upload to s3://%s/%s: %w", bucket, key, err)
	}

	return nil
}
