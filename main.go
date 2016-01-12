package main

import (
	"fmt"
	"os"
	"time"
)

const RetryTimeout = time.Second

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: retrysync <source> <destination>")
		os.Exit(1)
	}

	sourceStat := retryLStat(os.Args[1])
	if sourceStat.IsDir() {
		fmt.Println("Performing recursive copy...")
		if err := retryCopyDir(sourceStat, os.Args[1], os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "Fatal error:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Performing single file copy...")
		if err := retryCopyFile(sourceStat, os.Args[1], os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "Fatal error:", err)
			os.Exit(1)
		}
	}
}
