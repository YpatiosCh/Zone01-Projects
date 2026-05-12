package internal

import (
	"my-ls/internal/models"
	"os"
	"testing"
)

func TestCollectDirectoryEntries_All(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(dir+"/.hidden", []byte{}, 0644)
	os.WriteFile(dir+"/visible", []byte{}, 0644)

	entries, err := collectDirectoryEntries(dir, models.Flags{All: true})
	if err != nil {
		t.Fatal(err)
	}

	foundHidden := false
	for _, e := range entries {
		if e.Name == ".hidden" {
			foundHidden = true
		}
	}

	if !foundHidden {
		t.Fatal("expected hidden file with -a")
	}
}

func TestCollectDirectoryEntries_NoAll(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(dir+"/.hidden", []byte{}, 0644)

	entries, _ := collectDirectoryEntries(dir, models.Flags{})

	for _, e := range entries {
		if e.Name == ".hidden" {
			t.Fatal("hidden file should not appear without -a")
		}
	}
}
