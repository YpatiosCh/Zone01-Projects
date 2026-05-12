package writer

import (
	"bufio"
	"fmt"
	"os"
)

var Out = bufio.NewWriter(os.Stdout)

func Flush() {
	Out.Flush()
}

func PrintBlankLine() {
	fmt.Fprintln(Out)
}
