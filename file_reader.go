package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// retryFileReader is an io.Reader implementation which continually
// retries to read from a file until it succeeds, seeking as needed.
type retryFileReader struct {
	path   string
	offset int64
	handle *os.File
}

func newRetryFileReader(path string) *retryFileReader {
	return &retryFileReader{path: path}
}

func (r *retryFileReader) Read(output []byte) (n int, err error) {
	for {
		if r.handle == nil {
			var err error
			r.handle, err = os.Open(r.path)
			if err != nil {
				fmt.Println("Open("+r.path+") (offset", r.offset, "bytes) failed:", err)
				time.Sleep(RetryTimeout)
				continue
			}
			if r.offset != 0 {
				if _, err = r.handle.Seek(r.offset, 0); err != nil {
					fmt.Println("Seek(", r.offset, ") (in "+r.path+") failed:", err)
					time.Sleep(RetryTimeout)
					r.handle = nil
					continue
				}
			}
		}
		n, err = r.handle.Read(output)
		if err == nil || err == io.EOF {
			r.offset += int64(n)
			return
		}
		r.handle.Close()
		r.handle = nil

		fmt.Println("Read("+r.path+") (offset", r.offset, "bytes) failed:", err)
		time.Sleep(RetryTimeout)
	}
}

func (r *retryFileReader) Close() {
	if r.handle != nil {
		r.handle.Close()
	}
}
