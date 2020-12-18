package main

import (
	"context"
	"log"
)

func main() {
	interactor, err := build()
	if err != nil {
		log.Panicf("failed to build: %s", err.Error())
	}

	ctx := context.Background()
	if err := interactor.CreateRoleForCreatingThumbnails(ctx); err != nil {
		log.Panicf("failed to CreateRoleForCreatingThumbnails: %s", err.Error())
	}
}
