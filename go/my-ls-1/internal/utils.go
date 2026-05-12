package internal

import (
	"my-ls/internal/models"
	"os"
)

func isRealDir(e models.Entry) bool {
	if !e.Info.IsDir() {
		return false
	}
	// do not follow symlinks
	return e.Info.Mode()&os.ModeSymlink == 0
}
