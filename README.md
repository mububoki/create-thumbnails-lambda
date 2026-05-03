# create-thumbnails-lambda

This lambda function will be invoked once you put image to an Amazon S3 Bucket.

The function gets images from the bucket,
creates thumbnails,
and
put them to another bucket.

support: jpeg, gif, png

## image limits

Images that exceed any of these limits are rejected without producing a thumbnail.

| limit | value |
| :---: | :---: |
| file size | 50 MiB |
| width | 10000 px |
| height | 10000 px |

## environment variables

### lambda runtime

| name | default value | note |
| :---: | :---: | :--- |
| IMAGE_RATE | 0.5 | magnification rate |
| OBJECT_BUCKET_NAME_ORIGINAL | original.images.mububoki | Lambda will be invoked when you put images to this bucket. |
| OBJECT_BUCKET_NAME_THUMBNAIL | thumbnail.images.mububoki | Lambda puts thumbnails to this bucket. |

### setup CLI (cmd/setup)

| name | default value | note |
| :---: | :---: | :--- |
| AWS_REGION | (none) | required by AWS SDK |
| LAMBDA_FUNCTION_NAME | create-thumbnails-lambda | Lambda function name |
| LAMBDA_ROLE_NAME | create-thumbnails-lambda-role | IAM role name for the Lambda |
| LAMBDA_ZIP_PATH | bin/function.zip | path to the zip artifact for create/update |
| OBJECT_BUCKET_NAME_ORIGINAL | original.images.mububoki | source bucket name |
| OBJECT_BUCKET_NAME_THUMBNAIL | thumbnail.images.mububoki | destination bucket name |

## prerequisites

- Go (version specified in `go.mod`)
- `make`
- `zip`
- AWS credentials available via the standard SDK chain (`AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` env vars, `~/.aws/credentials`, etc.)

## make targets

### development

| target | description |
| :--- | :--- |
| `make build` | build for the host OS to `bin/app` |
| `make test` | run all tests |
| `make install-tools` | install `errcheck` |
| `make static-check` | `go vet` + `errcheck` |
| `make generate` | run `go generate` |

### Lambda artifact

| target | description |
| :--- | :--- |
| `make build-lambda-function` | cross-compile to `bin/bootstrap` (linux/arm64) |
| `make zip-lambda-function` | build then archive to `bin/function.zip` |

### AWS resource setup

| target | description |
| :--- | :--- |
| `make create-iam-role` | create the Lambda execution IAM role |
| `make create-s3-buckets` | create source/destination S3 buckets |
| `make delete-s3-buckets` | empty and delete source/destination S3 buckets |
| `make create-lambda-function` | create the Lambda function from `bin/function.zip` |
| `make update-lambda-function` | update the Lambda code from `bin/function.zip` |

## deployment

```sh
export AWS_REGION=ap-northeast-1
# AWS credentials must also be available via the SDK chain.

make create-iam-role
make create-s3-buckets
make zip-lambda-function
make create-lambda-function

# subsequent code updates:
make zip-lambda-function
make update-lambda-function
```

After the Lambda is created, configure the source bucket to invoke it on `s3:ObjectCreated:*` events via the AWS Console or another tool — the setup CLI does not configure the trigger.
