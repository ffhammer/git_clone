package commands

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/refs"
	"git_clone/gvc/settings"
	"git_clone/gvc/utils"

	"os"
	"path/filepath"
)

func InitGVC() string {
	// Check if the repository already exists
	err := utils.FindRepo()
	if err == nil {
		return fmt.Errorf("repository already exists at %s", utils.RepoDir).Error()
	}

	if err := os.Mkdir(config.OWN_FOLDER_NAME, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", config.OWN_FOLDER_NAME, err).Error()
	}

	objetsPath := filepath.Join(config.OWN_FOLDER_NAME, config.OBJECT_FOLDER)
	if err := os.Mkdir(objetsPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", objetsPath, err).Error()
	}

	if err := os.Mkdir(filepath.Join(config.OWN_FOLDER_NAME, config.RefsFolder), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", config.OWN_FOLDER_NAME, err).Error()
	}

	utils.RepoDir = config.OWN_FOLDER_NAME

	inital_metdata := refs.CurrentBranchPointer{ParentCommitHash: config.DOES_NOT_EXIST_HASH, BranchName: config.STARTING_BRANCH}
	err = refs.SaveCurrentPointer(inital_metdata)
	if err != nil {
		return fmt.Errorf("failed to save ref pointer: %w", err).Error()
	}

	setting := settings.Settings{LogLevel: settings.INFO, User: config.DOES_NOT_EXIST_HASH}
	if err := settings.SaveSettings(setting); err != nil {
		return fmt.Errorf("failed saving default settings: %w", err).Error()
	}

	return fmt.Sprintf("Initialized a new repository at %s", config.OWN_FOLDER_NAME)
}
