package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: setup <action>")
		fmt.Fprintln(os.Stderr, "actions: create-iam-role, create-s3-buckets, delete-s3-buckets")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "create-iam-role":
		err = createIAMRole()
	case "create-s3-buckets":
		err = createS3Buckets()
	case "delete-s3-buckets":
		err = deleteS3Buckets()
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
