package gvc

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func getPathToHashMap(treePath string) (map[string]string, error) {
	// Check if the file exists
	if _, err := os.Stat(treePath); err != nil {
		return nil, fmt.Errorf("could not find parent commit tree at %s: %w", treePath, err)
	}

	// Read the file data
	file, err := os.OpenFile(treePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", treePath, err)
	}

	// Parse the CSV data
	r := csv.NewReader(file)
	r.TrimLeadingSpace = true

	// Read the CSV header
	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}
	if len(headers) < 2 || headers[0] != "relPath" || headers[1] != "fileHash" {
		return nil, fmt.Errorf("invalid CSV format in file %s", treePath)
	}

	pathToHash := make(map[string]string)

	// Read and parse each line into the map
	for {
		record, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file, stop reading
			}
			return nil, fmt.Errorf("failed to read CSV line: %w", err)
		}

		if len(record) < 2 {
			continue // Skip lines that don't have enough columns
		}
		relPath := record[0]
		fileHash := record[1]
		pathToHash[relPath] = fileHash
	}

	return pathToHash, nil
}

func fileTableExistsOrCreate(tablePath string) error {
	if _, err := os.Stat(tablePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking file %s: %w", tablePath, err)
	}

	// create a csv with relPath, fileHash columns
	file, err := os.Create(tablePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", tablePath, err)
	}
	defer file.Close()

	// Initialize CSV writer and write header row
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written to the file

	headers := []string{"relPath", "fileHash"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers to file %s: %w", tablePath, err)
	}
	return nil

}

func addToSavedFilesTable(commitDir string, relPath string, fileHash string) error {

	tablePath := filepath.Join(commitDir, FILE_TABLE)

	// this should be a csv with relpath, Hash to store for a commit whats inside
	if err := fileTableExistsOrCreate(tablePath); err != nil {
		return fmt.Errorf("error for file table %s: %w", tablePath, err)
	}

	// Open the file for reading and writing
	file, err := os.OpenFile(tablePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file table %s: %w", tablePath, err)
	}
	defer file.Close()

	// Read existing records from the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil && err.Error() != "EOF" { // Allow empty files
		return fmt.Errorf("error reading file table: %w", err)
	}

	// Check if the file already exists in the records, and update if it does
	updated := false
	for i, record := range records {
		if len(record) >= 2 && record[0] == relPath {
			records[i][1] = fileHash
			updated = true
			break
		}
	}

	// If the file was not found, add a new record
	if !updated {
		records = append(records, []string{relPath, fileHash})
	}

	// Reset the file by seeking to the beginning and truncating it
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking to start of file: %w", err)
	}
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("error truncating file: %w", err)
	}

	// Write updated records back to the CSV file
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("error writing file table: %w", err)
	}
	writer.Flush()

	return nil

}
