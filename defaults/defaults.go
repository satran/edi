package defaults

import (
	"fmt"
	"os/user"
	"path/filepath"
)

func Root() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("can't find user home: %w", err)
	}
	return filepath.Join(usr.HomeDir, "notes"), nil
}
