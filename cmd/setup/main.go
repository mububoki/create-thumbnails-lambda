package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: setup <action>")
		fmt.Fprintln(os.Stderr, "actions: create-iam-role")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "create-iam-role":
		err = createIAMRole()
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
