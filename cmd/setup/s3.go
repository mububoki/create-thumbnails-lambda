package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
		input := &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		}

		// If region is not us-east-1, we must specify LocationConstraint
		if cfg.Region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(cfg.Region),
			}
		}

		if _, err := client.CreateBucket(ctx, input); err != nil {
			return fmt.Errorf("failed to CreateBucket %s: %w", bucket, err)
		}
		fmt.Printf("created bucket: %s in %s\n", bucket, cfg.Region)
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
		// Empty the bucket first
		if err := emptyBucket(ctx, client, bucket); err != nil {
			return fmt.Errorf("failed to empty bucket %s: %w", bucket, err)
		}

		if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(bucket),
		}); err != nil {
			return fmt.Errorf("failed to DeleteBucket %s: %w", bucket, err)
		}
		fmt.Printf("deleted bucket: %s\n", bucket)
	}

	return nil
}

func emptyBucket(ctx context.Context, client *s3.Client, bucket string) error {
	p := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return err
		}

		if len(page.Contents) == 0 {
			continue
		}

		var objects []types.ObjectIdentifier
		for _, obj := range page.Contents {
			objects = append(objects, types.ObjectIdentifier{Key: obj.Key})
		}

		if _, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &types.Delete{Objects: objects},
		}); err != nil {
			return err
		}
	}

	return nil
}
