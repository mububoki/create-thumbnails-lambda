package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

const (
	defaultLambdaZipPath    = "bin/function.zip"
	defaultLambdaTimeout    = 30
	defaultLambdaMemorySize = 256
)

func createLambdaFunction() error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	functionName := envOrDefault("LAMBDA_FUNCTION_NAME", defaultFunctionName)
	roleName := envOrDefault("LAMBDA_ROLE_NAME", defaultRoleName)
	zipPath := envOrDefault("LAMBDA_ZIP_PATH", defaultLambdaZipPath)
	bucketOriginal := envOrDefault("OBJECT_BUCKET_NAME_ORIGINAL", defaultBucketNameOriginal)
	bucketThumbnail := envOrDefault("OBJECT_BUCKET_NAME_THUMBNAIL", defaultBucketNameThumbnail)

	roleARN, err := getRoleARN(ctx, cfg, roleName)
	if err != nil {
		return err
	}

	zipBytes, err := os.ReadFile(zipPath)
	if err != nil {
		return fmt.Errorf("failed to read zip %s: %w", zipPath, err)
	}

	client := lambda.NewFromConfig(cfg)
	out, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName:  aws.String(functionName),
		Role:          aws.String(roleARN),
		Runtime:       lambdatypes.RuntimeProvidedal2023,
		Handler:       aws.String("bootstrap"),
		Code:          &lambdatypes.FunctionCode{ZipFile: zipBytes},
		Architectures: []lambdatypes.Architecture{lambdatypes.ArchitectureArm64},
		Timeout:       aws.Int32(defaultLambdaTimeout),
		MemorySize:    aws.Int32(defaultLambdaMemorySize),
		Environment: &lambdatypes.Environment{
			Variables: map[string]string{
				"OBJECT_BUCKET_NAME_ORIGINAL":  bucketOriginal,
				"OBJECT_BUCKET_NAME_THUMBNAIL": bucketThumbnail,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to CreateFunction: %w", err)
	}

	fmt.Printf("created function: %s\n", aws.ToString(out.FunctionArn))
	return nil
}

func updateLambdaFunction() error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to LoadDefaultConfig: %w", err)
	}

	functionName := envOrDefault("LAMBDA_FUNCTION_NAME", defaultFunctionName)
	zipPath := envOrDefault("LAMBDA_ZIP_PATH", defaultLambdaZipPath)

	zipBytes, err := os.ReadFile(zipPath)
	if err != nil {
		return fmt.Errorf("failed to read zip %s: %w", zipPath, err)
	}

	client := lambda.NewFromConfig(cfg)
	out, err := client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(functionName),
		ZipFile:      zipBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to UpdateFunctionCode: %w", err)
	}

	fmt.Printf("updated function: %s\n", aws.ToString(out.FunctionArn))
	return nil
}

func getRoleARN(ctx context.Context, cfg aws.Config, roleName string) (string, error) {
	client := iam.NewFromConfig(cfg)
	out, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to GetRole: %w", err)
	}
	return aws.ToString(out.Role.Arn), nil
}
