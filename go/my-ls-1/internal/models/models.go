package models

import (
	"os"
	"syscall"
)

type Flags struct {
	Long      bool // -l
	All       bool // -a
	Recursive bool // -R
	Dir       bool // -d
	One       bool // -1
	Reverse   bool // -r
	Time      bool // -t
}

type Entry struct {
	Name string
	Path string
	Info os.FileInfo
	Stat *syscall.Stat_t
}
