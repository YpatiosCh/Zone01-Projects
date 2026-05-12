package internal

import (
	"my-ls/internal/models"
	"os"
	"syscall"
)

func DispatchPath(path string, flags models.Flags) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	var st *syscall.Stat_t
	if s, ok := info.Sys().(*syscall.Stat_t); ok {
		st = s
	}
	entry := models.Entry{
		Name: path,
		Path: path,
		Info: info,
		Stat: st,
	}

	// -d means: treat directory like a file
	if info.IsDir() && !flags.Dir {
		return listDirectory(path, entry.Name, flags)
	}

	return listFile(entry, flags)
}
