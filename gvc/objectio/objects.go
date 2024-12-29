package objectio

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"io"
	"os"
	"path/filepath"
)

func SaveObject(fileHash string, reader io.Reader) error {

	subdir := filepath.Join(utils.RepoDir, config.OBJECT_FOLDER, fileHash[:2])

	if err := utils.MkdirIgnoreExists(subdir); err != nil {
		return fmt.Errorf("error creating subdir %s: %w", subdir, err)
	}

	// Define the full path where the compressed file will be saved
	objectFilePath := filepath.Join(subdir, fileHash[2:])
	if _, err := os.Stat(objectFilePath); err == nil {
		// File already exists, so return nil indicating no error (or you could return a specific error)
		return nil
	} else if !os.IsNotExist(err) {
		// If there's an error other than "file does not exist," return it
		return fmt.Errorf("error checking file %s: %w", objectFilePath, err)
	}
	destFile, err := os.Create(objectFilePath)
	if err != nil {
		return fmt.Errorf("error creating object file %s: %w", objectFilePath, err)
	}
	defer destFile.Close()

	// Compress the content and write it to the destination file
	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	// Copy the content from the source file to the gzip writer
	if _, err := io.Copy(gzipWriter, reader); err != nil {
		return fmt.Errorf("error compressing content: %w", err)
	}
	return nil
}

func LoadObject(fileHash string) ([]byte, error) {
	subdir := filepath.Join(utils.RepoDir, config.OBJECT_FOLDER, fileHash[:2])
	objectFilePath := filepath.Join(subdir, fileHash[2:])

	if _, err := os.Stat(objectFilePath); os.IsExist(err) {
		return nil, fmt.Errorf("Cant find object file %s: %w", objectFilePath, err)
	}

	file, err := os.Open(objectFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening object file %s: %w", objectFilePath, err)
	}
	defer file.Close()

	// Attempt to decompress the file content using gzip
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Read decompressed content
	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("error decompressing object file %s: %w", objectFilePath, err)
	}
	return content, nil
}

func AddFileToObjects(filename string, fileHash string) error {
	// Create the subdirectory using the first two characters of the hash

	// Open the source file to read its content
	sourceFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", filename, err)
	}
	defer sourceFile.Close()

	// Create the destination file (compressed)
	return SaveObject(fileHash, sourceFile)
}

func RetrieveFile(fileHash string) (string, error) {

	content, err := LoadObject(fileHash)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func SerializeObject[T any](object T) (io.Reader, error) {
	// Serialize the object to JSON
	data, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error serializing object: %w", err)
	}

	// Wrap the JSON data in a bytes.Reader to create an io.Reader
	return bytes.NewReader(data), nil
}

func DeserializeObject[T any](data []byte) (T, error) {
	var object T
	err := json.Unmarshal(data, &object)
	if err != nil {
		return object, fmt.Errorf("error deserializing object: %w", err)
	}
	return object, nil
}

func SaveJsonObject[T any](object T, fileHash string) error {
	reader, err := SerializeObject(object)
	if err != nil {
		return fmt.Errorf("error serializing tree: %w", err)
	}

	// Save the serialized tree as an object
	return SaveObject(fileHash, reader)
}

func LoadJsonObject[T any](fileHash string) (T, error) {
	// Load the serialized tree object
	data, err := LoadObject(fileHash)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error loading tree object: %w", err)
	}

	// Deserialize the data into a TreeMap
	return DeserializeObject[T](data)
}
