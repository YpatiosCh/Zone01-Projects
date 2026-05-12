package internal

import (
	"my-ls/internal/commands"
	"my-ls/internal/models"
	"my-ls/internal/writer"
	"my-ls/logger"
)

func listFile(entry models.Entry, flags models.Flags) error {

	if flags.Long {
		commands.PrintLong([]models.Entry{entry})
		return nil
	}
	if flags.One {
		commands.PrintEntry(entry, flags)
	} else {
		commands.PrintColumns([]models.Entry{entry})
	}
	return nil
}

func listDirectory(path, display string, flags models.Flags) error {
	entries, err := collectDirectoryEntries(path, flags)
	if err != nil {
		return err
	}

	sortEntries(entries, flags)

	// print this directory
	if flags.Long {
		commands.PrintLong(entries)
	} else {
		if flags.One {
			commands.PrintShort(entries)
		} else {
			commands.PrintColumns(entries)
		}
	}

	// recurse if -R
	if !flags.Recursive {
		return nil
	}

	for _, e := range entries {
		if !isRealDir(e) || e.Name == "." || e.Name == ".." {
			continue
		}

		writer.PrintBlankLine()
		nextDisplay := display + "/" + e.Name
		commands.PrintDirHeader(nextDisplay)

		if err := listDirectory(e.Path, nextDisplay, flags); err != nil {
			logger.PrintError(logger.FATAL, "", err)
		}
	}

	return nil
}
