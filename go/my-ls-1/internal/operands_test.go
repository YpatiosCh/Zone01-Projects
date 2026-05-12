package internal

import (
	"os"
	"testing"
)

func TestSplitOperands_FilesAndDirs(t *testing.T) {
	dir := t.TempDir()

	file := dir + "/file.txt"
	if err := os.WriteFile(file, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	files, dirs, errs := splitOperands([]string{file, dir})

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if len(dirs) != 1 {
		t.Fatalf("expected 1 dir, got %d", len(dirs))
	}
}

func TestSplitOperands_Error(t *testing.T) {
	_, _, errs := splitOperands([]string{"/does/not/exist"})

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}
