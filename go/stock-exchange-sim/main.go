package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/parser"
	"platform.zone01.gr/git/vstefano/stock-exchange-sim/scheduler"
)

// Test seller

func main() {
	filepath, timeout, err := handleArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	s, err := parser.NewParser().Run(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := scheduler.NewScheduler(s)

	err = sc.Optimize(context.Background(), timeout)
	if !errors.Is(err, scheduler.ErrNoMoreProcess) && err != nil {
		fmt.Println(err)
	}

}

func handleArgs() (filepath string, timeout time.Duration, err error) {
	argsLen := len(os.Args)

	if argsLen < 2 || argsLen > 3 {
		return "", 0, fmt.Errorf("usage: %s <file path> [OPTION] <duration in seconds>", os.Args[0])
	}

	filepath = os.Args[1]
	_, err = os.Stat(filepath)
	if err != nil {
		return "", 0, fmt.Errorf("file path error: %s, %w", filepath, err)
	}

	if argsLen > 2 {
		timeout, err = parseTime(os.Args[2])
		if err != nil {
			return "", 0, fmt.Errorf("timeout argument error: %s  %w", os.Args[2], err)
		}

	} else {
		timeout = time.Duration(float64(time.Second))
	}

	return filepath, timeout, nil
}

func parseTime(s string) (time.Duration, error) {

	s = strings.TrimSpace(s)

	// If the string ends with a known time unit, let time.ParseDuration handle it.
	if strings.HasSuffix(s, "ms") ||
		strings.HasSuffix(s, "s") {
		return time.ParseDuration(s)
	}

	// Otherwise, assume it's seconds.
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid time value: %q", s)
	}

	return time.Duration(v * float64(time.Second)), nil
}
