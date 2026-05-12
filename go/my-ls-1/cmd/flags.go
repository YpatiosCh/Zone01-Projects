package cmd

import (
	"fmt"
	"my-ls/internal/models"
	"my-ls/logger"
	"os"
)

func parseArgs(args []string) (models.Flags, []string, error) {
	var flags models.Flags
	var paths []string

	parsingFlags := true

	for _, arg := range args {
		if arg == "--help" {
			logger.Help()
			os.Exit(0)
		}
		if parsingFlags {
			if arg == "--" {
				parsingFlags = false
				continue
			}

			if len(arg) > 1 && arg[0] == '-' {
				for _, c := range arg[1:] {
					switch c {
					case 'l':
						flags.Long = true
					case 'a':
						flags.All = true
					case 'R':
						flags.Recursive = true
					case 'd':
						flags.Dir = true
					case '1':
						flags.One = true
					case 'r':
						flags.Reverse = true
					case 't':
						flags.Time = true

					default:
						return flags, nil, fmt.Errorf("invalid option -- %c", c)
					}
				}
				continue
			}
		}

		parsingFlags = false
		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	return flags, paths, nil
}
