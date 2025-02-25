package refs

import (
	"encoding/json"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/treediff"
	"git_clone/gvc/utils"
	"os"
	"path/filepath"
)

type MergeConflict struct {
	RelPath  string                `json:"relpath"`
	ActionA  treediff.ChangeAction `json:"ActionA"`
	NewHashA string                `json:"NewHashA"`
	ActionB  treediff.ChangeAction `json:"ActionB"`
	NewHashB string                `json:"NewHashB"`
}

type MergeMetaData struct {
	MERGE_HEAD    string `json:"MERGE_HEAD"`
	MERGE_MESSAGE string `json:"MERGE_MESSAGE"`
	Conflicts     []MergeConflict
}

func SaveMergeMetaData(data MergeMetaData) error {

	path := filepath.Join(utils.RepoDir, config.MERGE_INFO_PATH)

	byties, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error while saving merge meta data: json marshall failed %w", err)
	}

	err = os.WriteFile(path, byties, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error while saving merge meta data: writing to file '%s' failed %w", path, err)
	}

	return nil
}

func CurrentlyInMerge() (bool, error) {
	path := filepath.Join(utils.RepoDir, config.MERGE_INFO_PATH)

	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, nil
}

func GetMergeMetaData() (MergeMetaData, error) {
	path := filepath.Join(utils.RepoDir, config.MERGE_INFO_PATH)
	byites, err := os.ReadFile(path)
	if err != nil {
		return MergeMetaData{}, fmt.Errorf("error while loading merge meta data from '%s': %w", path, err)
	}

	data := MergeMetaData{}
	err = json.Unmarshal(byites, &data)
	if err != nil {
		return MergeMetaData{}, fmt.Errorf("error while loading merge meta data from '%s': json unmarshal %w", path, err)
	}

	return data, nil
}

func DelMergeMetaMData() error {
	path := filepath.Join(utils.RepoDir, config.MERGE_INFO_PATH)
	err := os.Remove(path)
	if err != nil {
		return nil
	}

	return fmt.Errorf("could not delete merge meta data '%s' %w", path, err)
}
