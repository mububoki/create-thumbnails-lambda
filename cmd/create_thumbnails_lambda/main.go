package main

import (
	"log"
)

func main() {
	handler, err := build()
	if err != nil {
		log.Panicf("failed to build: %s", err.Error())
	}
	handler.CreateThumbnail()
}
