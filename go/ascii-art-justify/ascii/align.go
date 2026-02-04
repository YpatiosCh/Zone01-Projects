package ascii

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// create justify func

func GetTerminalWidth() (int, error) {
	var winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&winsize)))
	if err != 0 {
		return 0, err
	}

	return int(winsize.Col), nil
}

func AlignCenter(lines []string, width, length int) {
	for _, line := range lines {
		centerLine(line, width, length)
	}
}

func centerLine(line string, termWidth, length int) {
	padding := (termWidth - length) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Println(strings.Repeat(" ", padding) + line)
}

func rightLine(line string, termWidth, length int) {
	padding := (termWidth - length)
	if padding < 0 {
		padding = 0
	}
	fmt.Println(strings.Repeat(" ", padding) + line)
}

func AlignRight(lines []string, width, length int) {
	for _, line := range lines {
		rightLine(line, width, length)
	}
}

func Justify(art []string, length int, word string) {
	width, _ := GetTerminalWidth()
	toFindIndex := strings.Split(word, " ")
	index := len(toFindIndex)
	padding := (width - length) / (index - 1)
	var pad string
	for i := 1; i < padding; i++ {
		pad += " "
	}
	for i := range art {
		art[i] += pad
	}
}

// Validate the color flag to ensure it uses "=" syntax
func ValidateAlignFlag() error {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--align") && !strings.Contains(arg, "=") {
			return fmt.Errorf("\ninvalid syntax for --align flag. please use '--align=<type>' format")
		}
	}
	return nil
}

func ValidateAlignment(colorFlag string) error {
	if colorFlag != "center" && colorFlag != "right" && colorFlag != "justify" && colorFlag != "left" && colorFlag != "" {
		return fmt.Errorf("\ninvalid alignment! You can choose between <left> <right> <center> <justify>")
	}
	return nil
}
