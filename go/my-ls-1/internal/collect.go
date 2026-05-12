package internal

import (
	"my-ls/internal/models"
	"os"
	"path/filepath"
	"syscall"
)

var entries []models.Entry

func collectDirectoryEntries(path string, flags models.Flags) ([]models.Entry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []models.Entry

	if flags.All {
		// .
		if info, err := os.Lstat(path); err == nil {
			entries = append(entries, models.Entry{
				Name: ".",
				Path: path,
				Info: info,
			})
		}

		// ..
		parent := filepath.Dir(path)
		if info, err := os.Lstat(parent); err == nil {
			entries = append(entries, models.Entry{
				Name: "..",
				Path: parent,
				Info: info,
			})
		}
	}

	for _, de := range dirEntries {
		name := de.Name()

		// -a: include dotfiles
		if !flags.All && len(name) > 0 && name[0] == '.' {
			continue
		}

		fullPath := filepath.Join(path, name)

		info, err := os.Lstat(fullPath)
		if err != nil {
			// ls skips unreadable entries but continues
			continue
		}

		var st *syscall.Stat_t
		if s, ok := info.Sys().(*syscall.Stat_t); ok {
			st = s
		}

		entries = append(entries, models.Entry{
			Name: name,
			Path: fullPath,
			Info: info,
			Stat: st,
		})
	}

	return entries, nil
}
