package main

import (
	"fmt"
	"os"
	"time"
)

func retryStat(path string) os.FileInfo {
	for {
		if res, err := os.Stat(path); err == nil {
			return res
		} else {
			fmt.Println("Stat("+path+") failed:", err)
		}
		time.Sleep(RetryTimeout)
	}
}

func retryListDir(path string) (listing []os.FileInfo) {
	for {
		var err error
		var f *os.File
		if f, err = os.Open(path); err == nil {
			listing, err = f.Readdir(0)
			f.Close()
			if err == nil {
				return
			}
		}
		fmt.Println("Readdir("+path+") failed:", err)
		time.Sleep(RetryTimeout)
	}
}
