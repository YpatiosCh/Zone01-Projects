package internal

import (
	"my-ls/internal/models"
	"os"
	"syscall"
)

type AccessError struct {
	Path string
	Err  error
}

func splitOperands(paths []string) (files []models.Entry, dirs []models.Entry, errs []AccessError) {
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			errs = append(errs, AccessError{Path: path, Err: err})
			continue
		}

		var st *syscall.Stat_t
		if s, ok := info.Sys().(*syscall.Stat_t); ok {
			st = s
		}
		e := models.Entry{
			Name: path,
			Path: path,
			Info: info,
			Stat: st,
		}

		// -d is handled outside (Dispatch/Execute logic), so here we classify raw operands.
		if info.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	return
}
