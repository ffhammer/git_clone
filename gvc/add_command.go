package gvc

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addFileToObjects(repoDir string, filename string, fileHash string) error {
	// Create the subdirectory using the first two characters of the hash
	subdir := filepath.Join(repoDir, OBJECT_FOLDER, fileHash[:2])

	if err := mkdirIgnoreExists(subdir); err != nil {
		return fmt.Errorf("error creating subdir %s: %w", subdir, err)
	}

	// Define the full path where the compressed file will be saved
	objectFilePath := filepath.Join(subdir, fileHash)

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

func AddFile(repoDir string, filePath string, force bool) error {

	// place we store things we added but have yet to commit
	next_commit := filepath.Join(repoDir, NEXT_COMMIT)

	mkdirIgnoreExists(next_commit)

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("can't find file %s: %w", filePath, err)
	}

	relPath, err := makePathRelativeToRepo(repoDir, filePath)
	if err != nil {
		return err
	}

	if isInIgnoreFile(relPath, repoDir) && !force {
		return fmt.Errorf("file %s is ignored. Use add -f to force it", filePath)
	}

	fileHash, err := getFileSHA256(filePath)
	if err != nil {
		return fmt.Errorf("can't add file %s to objects because %w", filePath, err)

	}

	err = addFileToObjects(repoDir, filePath, fileHash)
	if err != nil {
		return fmt.Errorf("can't add file %s to objects because %w", filePath, err)
	}

	if err := addToSavedFilesTable(next_commit, relPath, fileHash); err != nil {
		return err
	}

	return nil
}
