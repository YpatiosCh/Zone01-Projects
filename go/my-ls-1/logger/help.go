package logger

import (
	"fmt"

	"github.com/Des1red/clihelp"
)

func Help() {
	fmt.Println("Usage: my-ls [OPTION]... [FILE]...")
	fmt.Println()
	fmt.Println("List information about the FILEs (the current directory by default).")
	fmt.Println()
	fmt.Println("Flags:")

	clihelp.Print(
		clihelp.F("-l", "", "use a long listing format"),
		clihelp.F("-a", "", "do not ignore entries starting with ."),
		clihelp.F("-R", "", "list subdirectories recursively"),
		clihelp.F("-d", "", "list directories themselves, not their contents"),
		clihelp.F("-1", "", "list one file per line"),
		clihelp.F("-r", "", "reverse order while sorting"),
		clihelp.F("-t", "", "sort by modification time, newest first"),
		clihelp.F("--help", "", "display this help and exit"),
	)
}
