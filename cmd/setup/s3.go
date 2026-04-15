package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func createS3Buckets() error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	buckets := []string{
		envOrDefault("OBJECT_BUCKET_NAME_ORIGINAL", defaultBucketNameOriginal),
		envOrDefault("OBJECT_BUCKET_NAME_THUMBNAIL", defaultBucketNameThumbnail),
	}

	for _, bucket := range buckets {
		if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		}); err != nil {
			return fmt.Errorf("failed to CreateBucket %s: %w", bucket, err)
		}
		fmt.Printf("created bucket: %s\n", bucket)
	}

	return nil
}

func deleteS3Buckets() error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	buckets := []string{
		envOrDefault("OBJECT_BUCKET_NAME_ORIGINAL", defaultBucketNameOriginal),
		envOrDefault("OBJECT_BUCKET_NAME_THUMBNAIL", defaultBucketNameThumbnail),
	}

	for _, bucket := range buckets {
		if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(bucket),
		}); err != nil {
			return fmt.Errorf("failed to DeleteBucket %s: %w", bucket, err)
		}
		fmt.Printf("deleted bucket: %s\n", bucket)
	}

	return nil
}
