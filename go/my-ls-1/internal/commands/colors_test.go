package commands

import (
	"my-ls/internal/models"
	"os"
	"strings"
	"testing"
	"time"
)

type fakeDirInfo struct{}

func (f fakeDirInfo) Name() string       { return "dir" }
func (f fakeDirInfo) Size() int64        { return 0 }
func (f fakeDirInfo) Mode() os.FileMode  { return os.ModeDir }
func (f fakeDirInfo) ModTime() time.Time { return time.Now() }
func (f fakeDirInfo) IsDir() bool        { return true }
func (f fakeDirInfo) Sys() any           { return nil }

func TestFormatColoredName_Disabled(t *testing.T) {
	enableColor = false

	e := models.Entry{Name: "dir"}
	result := formatColoredName(e)

	if strings.Contains(result, "\033[") {
		t.Fatal("unexpected ANSI color when disabled")
	}
}

func TestFormatColoredName_Enabled(t *testing.T) {
	enableColor = true

	e := models.Entry{Name: "dir", Info: fakeDirInfo{}}
	result := formatColoredName(e)

	if !strings.Contains(result, "\033[34m") {
		t.Fatal("expected blue color for directory")
	}
}
