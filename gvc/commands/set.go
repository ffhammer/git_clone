package commands

import (
	"flag"
	"fmt"
	"git_clone/gvc/settings"
)

func listSettings() (string, error) {
	cfg, err := settings.LoadSettings()
	if err != nil {
		return "", err
	}
	return cfg.List() // Calls the exported List() method on settings
}

func setSetting(input string) error {
	cfg, err := settings.LoadSettings()
	if err != nil {
		return err
	}
	// Use a pointer receiver method to update the setting.
	if err := cfg.Set(input); err != nil {
		return err
	}
	return settings.SaveSettings(cfg)
}

func SettingsCommand(inputArgs []string) string {
	flagset := flag.NewFlagSet("settings", flag.ExitOnError)
	isSet := flagset.Bool("set", false, "set new value. Must provide key=value")
	isList := flagset.Bool("list", false, "list settings")

	if err := flagset.Parse(inputArgs); err != nil {
		return fmt.Errorf("error parsing args: %w", err).Error()
	}

	// Ensure that exactly one option is provided.
	if (*isList && *isSet) || (!*isList && !*isSet) {
		return "Choose exactly one option: either --list or --set"
	}

	if *isList {
		output, err := listSettings()
		if err != nil {
			return err.Error()
		}
		return output
	}

	// We know --set is true.
	if len(flagset.Args()) != 1 {
		return "Need one argument in key=value format"
	}

	if err := setSetting(flagset.Arg(0)); err != nil {
		return err.Error()
	}
	return "Setting updated successfully"
}
