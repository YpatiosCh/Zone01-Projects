package internal

import (
	"my-ls/internal/models"
	"testing"
)

func TestListDirectory_Error(t *testing.T) {
	err := listDirectory("/does/not/exist", ".", models.Flags{})
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}
