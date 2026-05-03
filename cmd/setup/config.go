package main

import "os"

const (
	defaultBucketNameOriginal   = "original.images.mububoki"
	defaultBucketNameThumbnail  = "thumbnail.images.mububoki"
	defaultFunctionName         = "create-thumbnails-lambda"
	defaultRoleName             = "create-thumbnails-lambda-role"
	lambdaBasicExecutionRoleARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
	s3PolicyName                = "create-thumbnails-s3-access"
)

type trustPolicy struct {
	Version   string           `json:"Version"`
	Statement []trustStatement `json:"Statement"`
}

type trustStatement struct {
	Effect    string         `json:"Effect"`
	Principal trustPrincipal `json:"Principal"`
	Action    string         `json:"Action"`
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
