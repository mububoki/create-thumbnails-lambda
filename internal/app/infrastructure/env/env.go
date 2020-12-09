package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Image ImageEnv
var Object ObjectEnv

type ImageEnv struct {
	Rate float64 `envconfig:"RATE" default:"0.5"`
}

type ObjectEnv struct {
	BucketNameOriginal  string `envconfig:"BUCKET_NAME_ORIGINAL" default:"original.images.mububoki"`
	BucketNameThumbnail string `envconfig:"BUCKET_NAME_THUMBNAIL" default:"thumbnail.images.mububoki"`
}

func init() {
	if err := envconfig.Process("IMAGE", &Image); err != nil {
		log.Panicf("failed to Process: %s", err.Error())
	}
	if err := envconfig.Process("OBJECT", &Object); err != nil {
		log.Panicf("failed to Process: %s", err.Error())
	}
}
