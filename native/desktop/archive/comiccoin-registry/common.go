package main

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	// Check if the source file exists
	if _, err := os.Stat(src); err != nil {
		return err
	}

	// Check if the destination file is writable
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination file already exists")
	}

	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Open the destination file for writing
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
