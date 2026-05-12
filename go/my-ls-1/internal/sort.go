package internal

import (
	"my-ls/internal/models"
	"sort"
	"strings"
)

func sortEntries(entries []models.Entry, flags models.Flags) {
	sort.Slice(entries, func(i, j int) bool {
		a := entries[i]
		b := entries[j]

		// "." first
		if a.Name == "." {
			return b.Name != "."
		}
		if b.Name == "." {
			return false
		}

		// ".." second
		if a.Name == ".." {
			return b.Name != ".."
		}
		if b.Name == ".." {
			return false
		}

		var less bool

		// -t: sort by mod time (newest first)
		if flags.Time {
			at := a.Info.ModTime()
			bt := b.Info.ModTime()
			if !at.Equal(bt) {
				less = at.After(bt)
			} else {
				less = nameLess(a.Name, b.Name)
			}
		} else {
			less = nameLess(a.Name, b.Name)
		}

		// -r: reverse result
		if flags.Reverse {
			return !less
		}
		return less
	})
}

func nameLess(a, b string) bool {
	ka := strings.TrimPrefix(a, ".")
	kb := strings.TrimPrefix(b, ".")

	la := strings.ToLower(ka)
	lb := strings.ToLower(kb)

	if la != lb {
		return la < lb
	}
	return a < b
}
