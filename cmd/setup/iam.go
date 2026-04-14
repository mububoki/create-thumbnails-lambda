package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

const (
	defaultRoleName           = "create-thumbnails-lambda-role"
	defaultBucketNameOriginal = "original.images.mububoki"
	defaultBucketNameThumbnail = "thumbnail.images.mububoki"
	lambdaBasicExecutionRoleARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
	s3PolicyName               = "create-thumbnails-s3-access"
)

type trustPolicy struct {
	Version   string            `json:"Version"`
	Statement []trustStatement  `json:"Statement"`
}

type trustStatement struct {
	Effect    string          `json:"Effect"`
	Principal trustPrincipal  `json:"Principal"`
	Action    string          `json:"Action"`
}

type trustPrincipal struct {
	Service string `json:"Service"`
}

type s3Policy struct {
	Version   string        `json:"Version"`
	Statement []s3Statement `json:"Statement"`
}

type s3Statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

func envOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func createIAMRole() error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	client := iam.NewFromConfig(cfg)
	roleName := envOrDefault("LAMBDA_ROLE_NAME", defaultRoleName)
	bucketOriginal := envOrDefault("OBJECT_BUCKET_NAME_ORIGINAL", defaultBucketNameOriginal)
	bucketThumbnail := envOrDefault("OBJECT_BUCKET_NAME_THUMBNAIL", defaultBucketNameThumbnail)

	// Create trust policy for Lambda
	tp := trustPolicy{
		Version: "2012-10-17",
		Statement: []trustStatement{
			{
				Effect:    "Allow",
				Principal: trustPrincipal{Service: "lambda.amazonaws.com"},
				Action:    "sts:AssumeRole",
			},
		},
	}
	trustPolicyJSON, err := json.Marshal(tp)
	if err != nil {
		return fmt.Errorf("failed to marshal trust policy: %w", err)
	}

	// Create IAM role
	createRoleOutput, err := client.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(string(trustPolicyJSON)),
		Tags: []types.Tag{
			{Key: aws.String("Project"), Value: aws.String("create-thumbnails-lambda")},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to CreateRole: %w", err)
	}

	fmt.Printf("created role: %s\n", *createRoleOutput.Role.Arn)

	// Attach AWSLambdaBasicExecutionRole managed policy
	if _, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(lambdaBasicExecutionRoleARN),
	}); err != nil {
		return fmt.Errorf("failed to AttachRolePolicy: %w", err)
	}

	fmt.Println("attached AWSLambdaBasicExecutionRole")

	// Create and attach inline policy for S3 access
	sp := s3Policy{
		Version: "2012-10-17",
		Statement: []s3Statement{
			{
				Effect:   "Allow",
				Action:   []string{"s3:GetObject"},
				Resource: fmt.Sprintf("arn:aws:s3:::%s/*", bucketOriginal),
			},
			{
				Effect:   "Allow",
				Action:   []string{"s3:PutObject"},
				Resource: fmt.Sprintf("arn:aws:s3:::%s/*", bucketThumbnail),
			},
		},
	}
	s3PolicyJSON, err := json.Marshal(sp)
	if err != nil {
		return fmt.Errorf("failed to marshal s3 policy: %w", err)
	}

	if _, err := client.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(s3PolicyName),
		PolicyDocument: aws.String(string(s3PolicyJSON)),
	}); err != nil {
		return fmt.Errorf("failed to PutRolePolicy: %w", err)
	}

	fmt.Println("attached S3 access policy")
	fmt.Printf("role ARN: %s\n", *createRoleOutput.Role.Arn)

	return nil
}
