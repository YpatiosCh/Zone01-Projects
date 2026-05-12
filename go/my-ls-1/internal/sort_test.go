package internal

import (
	"my-ls/internal/models"
	"os"
	"testing"
	"time"
)

type fakeFileInfo struct {
	name string
	mod  time.Time
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f fakeFileInfo) ModTime() time.Time { return f.mod }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }

func entry(name string, mtime time.Time) models.Entry {
	return models.Entry{
		Name: name,
		Info: fakeFileInfo{name: name, mod: mtime},
	}
}

func TestSortEntries_Default(t *testing.T) {
	entries := []models.Entry{
		entry("b", time.Now()),
		entry("a", time.Now()),
		entry("c", time.Now()),
	}

	sortEntries(entries, models.Flags{})

	if entries[0].Name != "a" || entries[1].Name != "b" {
		t.Fatalf("unexpected order: %v", entries)
	}
}

func TestSortEntries_Reverse(t *testing.T) {
	entries := []models.Entry{
		entry("a", time.Now()),
		entry("b", time.Now()),
	}

	sortEntries(entries, models.Flags{Reverse: true})

	if entries[0].Name != "b" {
		t.Fatalf("expected b first, got %s", entries[0].Name)
	}
}

func TestSortEntries_Time(t *testing.T) {
	old := time.Now().Add(-time.Hour)
	newer := time.Now()

	entries := []models.Entry{
		entry("old", old),
		entry("new", newer),
	}

	sortEntries(entries, models.Flags{Time: true})

	if entries[0].Name != "new" {
		t.Fatalf("expected newest first, got %s", entries[0].Name)
	}
}
