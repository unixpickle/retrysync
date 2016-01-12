package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func retryCopyDir(sourceInfo os.FileInfo, source, destination string) error {
	if stat, err := os.Stat(destination); err == nil {
		if !stat.IsDir() {
			return errors.New("target exists but is not a directory: " + destination)
		}
		fmt.Println("Using existing destination:", destination)
	} else {
		if err := os.Mkdir(destination, sourceInfo.Mode()&os.ModePerm); err != nil {
			return err
		}
	}

	listing := retryListDir(source)
	for _, info := range listing {
		newSource := filepath.Join(source, info.Name())
		newDest := filepath.Join(destination, info.Name())
		if info.IsDir() {
			if err := retryCopyDir(info, newSource, newDest); err != nil {
				return err
			}
		} else {
			if err := retryCopyFile(info, newSource, newDest); err != nil {
				return err
			}
		}
	}

	// Change the times after copying the contents to avoid updating the modification time.
	if err := os.Chtimes(destination, sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
		return err
	}

	return nil
}

func retryCopyFile(sourceInfo os.FileInfo, source, destination string) error {
	if destInfo, err := os.Lstat(destination); err == nil {
		if destInfo.IsDir() {
			return errors.New("target exists but is a directory: " + destination)
		}
		if (destInfo.Mode() & os.ModeSymlink) != (sourceInfo.Mode() & os.ModeSymlink) {
			return errors.New("target exists but symlink-ness differs: " + destination)
		}
		if destInfo.Size() == sourceInfo.Size() {
			fmt.Println("Skipping file:", destination)
			return nil
		}
		fmt.Println("Overwriting file:", destination)
		os.Remove(destination)
	}

	if (sourceInfo.Mode() & os.ModeSymlink) != 0 {
		fmt.Println("Copying link:", source)
		linkDest := retryReadlink(source)
		return os.Symlink(linkDest, destination)
	}

	fmt.Println("Copying file:", source)

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	reader := newRetryFileReader(source)
	_, err = io.Copy(output, reader)
	reader.Close()
	output.Close()

	if err != nil {
		return err
	}
	return os.Chtimes(destination, sourceInfo.ModTime(), sourceInfo.ModTime())
}
