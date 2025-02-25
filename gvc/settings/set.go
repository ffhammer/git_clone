package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/utils"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type LogLevel string

const (
	INFO      LogLevel = "INFO"
	DEBUGGING LogLevel = "DEBUGGING"
	WARNING   LogLevel = "WARNING"
	ERROR     LogLevel = "ERROR"
)

// Settings represents user configuration
type Settings struct {
	User     string   `json:"User"`
	LogLevel LogLevel `json:"LogLevel"`
}

// Set updates a setting field based on an input string in "key=value" format.
// It requires that the field is of type string.
func (setting *Settings) Set(input string) error {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return errors.New("error settings a 'Settings' value: expected key=value")
	}
	key := parts[0]
	newValue := parts[1]

	// Get a settable reflection value for the underlying struct.
	v := reflect.ValueOf(setting).Elem()
	t := v.Type()

	// Iterate over the struct fields.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name != key {
			continue
		}

		f := v.Field(i)
		if !f.CanSet() {
			return fmt.Errorf("cannot set field %s", key)
		}

		// Ensure the field is of string type.
		if f.Kind() != reflect.String {
			return errors.New("error settings a 'Settings' value: \nlogic error: settings must be strings at the moment")
		}

		// Set the field to the new string value.
		f.SetString(newValue)
		return nil
	}

	return fmt.Errorf("error settings a 'Settings' value: no such field: %s", key)
}

func (setting Settings) List() (string, error) {
	var builder strings.Builder
	v := reflect.ValueOf(setting)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		f := v.Field(i)
		if f.Kind() != reflect.String {
			return "", errors.New("logic error while listint setting: settings must be strings at the moment")
		}
		builder.WriteString(fmt.Sprintf("%s = '%s'\n", field.Name, f.String()))
	}
	return builder.String(), nil
}

// LoadSettings loads the settings from a JSON file.
// If the file does not exist, it returns an empty settings struct.
func LoadSettings() (Settings, error) {
	path := filepath.Join(utils.RepoDir, config.SETTINGS_PATH)

	// Open the settings file for reading
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		// If the file does not exist, return an empty settings struct with defaults
		return Settings{}, nil
	}

	if err != nil {
		return Settings{}, fmt.Errorf("error while loading settings: %w", err)
	}
	defer file.Close()

	// Read the entire file content
	data, err := io.ReadAll(file)
	if err != nil {
		return Settings{}, fmt.Errorf("error reading settings file: %w", err)
	}

	// Deserialize JSON data into Settings struct
	var cfgs Settings
	if err := json.Unmarshal(data, &cfgs); err != nil {
		return Settings{}, fmt.Errorf("error parsing settings file: %w", err)
	}

	return cfgs, nil
}

// SaveSettings serializes and writes the settings to a JSON file.
func SaveSettings(cfgs Settings) error {
	data, err := json.MarshalIndent(cfgs, "", "  ") // Pretty print JSON
	if err != nil {
		return fmt.Errorf("error serializing settings: %w", err)
	}

	path := filepath.Join(utils.RepoDir, config.SETTINGS_PATH)

	// Ensure parent directories exist before writing
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating settings directory: %w", err)
	}

	// Open the settings file for writing
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening settings file at %s: %w", path, err)
	}
	defer file.Close()

	// Write JSON data to the file
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("error writing to settings file at %s: %w", path, err)
	}

	return nil
}
