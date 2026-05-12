package cmd

import (
	"my-ls/internal"
	"my-ls/logger"
	"os"
)

func Execute() {
	flags, paths, err := parseArgs(os.Args[1:])
	if err != nil {
		logger.PrintError(logger.FATAL, "", err)
	}

	internal.Run(paths, flags)

	logger.ExitStatus()
}
