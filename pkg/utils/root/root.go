package root

import (
	"strconv"

	"github.com/pkg/errors"
	cmd "nenvoy.com/pkg/utils/cmd"
)

// AsRoot - Check running as root user (Linux, macOS)
func AsRoot() (bool, error) {
	// Check the UID for the user who ran the application
	stdout, err := cmd.Run("id", "-u")
	if err != nil {
		return false, errors.Wrap(err, "failed on AsRoot (1):")
	}

	//
	i, err := strconv.Atoi(string(stdout[:len(stdout)-1]))
	if err != nil {
		return false, errors.Wrap(err, "failed on AsRoot (2):")
	}

	if i == 0 {
		return true, nil
	}

	return false, nil
}
