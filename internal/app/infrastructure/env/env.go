package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Image ImageEnv
var Object ObjectEnv
var AuthNZ AuthNZEnv

type ImageEnv struct {
	Rate float64 `envconfig:"RATE" default:"0.5"`
}

type ObjectEnv struct {
	BucketNameOriginal  string `envconfig:"BUCKET_NAME_ORIGINAL" default:"original.images.mububoki"`
	BucketNameThumbnail string `envconfig:"BUCKET_NAME_THUMBNAIL" default:"thumbnail.images.mububoki"`
}

type AuthNZEnv struct {
	RoleName string   `envconfig:"ROLE_NAME" default:"lambda-role-create-thumbnails"`
	Actions  []string `envconfig:"ACTIONS" default:"[s3:PutObject, s3;GetObject]"`
}

func init() {
	if err := envconfig.Process("IMAGE", &Image); err != nil {
		log.Panicf("failed to Process: %s", err.Error())
	}
	if err := envconfig.Process("OBJECT", &Object); err != nil {
		log.Panicf("failed to Process: %s", err.Error())
	}
	if err := envconfig.Process("AUTH_NZ", &AuthNZ); err != nil {
		log.Panicf("failed to Process: %s", err.Error())
	}
}
