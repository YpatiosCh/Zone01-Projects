package commands

import (
	"fmt"
	"my-ls/internal/models"
	"my-ls/internal/writer"
	"os"
	"syscall"
)

var userCache = make(map[uint32]string)
var groupCache = make(map[uint32]string)

func PrintLong(entries []models.Entry) {
	if len(entries) > 1 {
		fmt.Fprintf(writer.Out, "total %d\n", totalBlocks(entries))
	}
	var maxLinks, maxSize int
	var maxUser, maxGroup int

	for _, e := range entries {
		st, ok := statOK(e)
		if !ok {
			continue
		}

		links := int(st.Nlink)
		size := int(e.Info.Size())

		u := username(st.Uid)
		g := groupname(st.Gid)

		if digits(links) > maxLinks {
			maxLinks = digits(links)
		}
		if digits(size) > maxSize {
			maxSize = digits(size)
		}
		if len(u) > maxUser {
			maxUser = len(u)
		}
		if len(g) > maxGroup {
			maxGroup = len(g)
		}
	}

	for _, e := range entries {
		st, ok := statOK(e)
		if !ok {
			continue
		}

		fmt.Fprintf(writer.Out,
			"%s %*d %-*s %-*s %*d %s %s\n",
			e.Info.Mode().String(),
			maxLinks, st.Nlink,
			maxUser, username(st.Uid),
			maxGroup, groupname(st.Gid),
			maxSize, e.Info.Size(),
			formatTime(e.Info.ModTime()),
			formatName(e),
		)
	}

}

func statOK(e models.Entry) (*syscall.Stat_t, bool) {
	return e.Stat, e.Stat != nil
}

func totalBlocks(entries []models.Entry) int64 {
	var total512 int64
	for _, e := range entries {
		st, ok := statOK(e)
		if !ok {
			continue
		}
		total512 += st.Blocks // 512-byte blocks on Linux
	}

	// Convert 512B blocks to 1K blocks (round up)
	return (total512 + 1) / 2
}

func formatName(e models.Entry) string {
	name := e.Name

	if enableColor && e.Info.IsDir() {
		name = colorBlue + name + colorReset
	}

	if e.Info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(e.Path); err == nil {
			return name + " -> " + target
		}
	}
	return name
}
