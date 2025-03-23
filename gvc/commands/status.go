package commands

import (
	"fmt"
	"git_clone/gvc/status"
)

func Status() string {
	output, err := status.Status()
	if err != nil {
		return fmt.Errorf("status failed because %w", err).Error()
	}
	return output
}
