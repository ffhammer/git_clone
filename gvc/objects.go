package gvc

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addFileToObjects(filename string, fileHash string) error {
	// Create the subdirectory using the first two characters of the hash
	subdir := filepath.Join(repoDir, OBJECT_FOLDER, fileHash[:2])

	if err := mkdirIgnoreExists(subdir); err != nil {
		return fmt.Errorf("error creating subdir %s: %w", subdir, err)
	}

	// Define the full path where the compressed file will be saved
	objectFilePath := filepath.Join(subdir, fileHash[2:])

	// Check if the file already exists
	if _, err := os.Stat(objectFilePath); err == nil {
		// File already exists, so return nil indicating no error (or you could return a specific error)
		return nil
	} else if !os.IsNotExist(err) {
		// If there's an error other than "file does not exist," return it
		return fmt.Errorf("error checking file %s: %w", objectFilePath, err)
	}

	// Open the source file to read its content
	sourceFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", filename, err)
	}
	defer sourceFile.Close()

	// Create the destination file (compressed)
	destFile, err := os.Create(objectFilePath)
	if err != nil {
		return fmt.Errorf("error creating object file %s: %w", objectFilePath, err)
	}
	defer destFile.Close()

	// Compress the content and write it to the destination file
	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	// Copy the content from the source file to the gzip writer
	if _, err := io.Copy(gzipWriter, sourceFile); err != nil {
		return fmt.Errorf("error compressing content: %w", err)
	}

	return nil
}

func retrieveObject(fileHash string) (string, error) {
	subdir := filepath.Join(repoDir, OBJECT_FOLDER, fileHash[:2])
	objectFilePath := filepath.Join(subdir, fileHash[2:])

	if _, err := os.Stat(objectFilePath); os.IsExist(err) {
		return "", fmt.Errorf("Cant find object file %s: %w", objectFilePath, err)
	}

	file, err := os.Open(objectFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening object file %s: %w", objectFilePath, err)
	}
	defer file.Close()

	// Attempt to decompress the file content using gzip
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		// If it's not a gzip file, we assume it’s plain text, so we rewind and read normally
		file.Seek(0, 0)
		content, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("error reading object file %s: %w", objectFilePath, err)
		}
		return string(content), nil
	}
	defer gzipReader.Close()

	// Read decompressed content
	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", fmt.Errorf("error decompressing object file %s: %w", objectFilePath, err)
	}

	return string(content), nil
}
