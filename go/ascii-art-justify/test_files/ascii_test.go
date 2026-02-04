package test_files

import (
	"ascii-art-justify/ascii"
	"testing"
)

func TestIsSubstring(t *testing.T) {
	if !ascii.IsSubstring("hello world", "world") {
		t.Error("expected 'world' to be a substring of 'hello world'")
	}
	if ascii.IsSubstring("hello world", "test") {
		t.Error("expected 'test' not to be a substring of 'hello world'")
	}
}

func TestIsFont(t *testing.T) {
	if !ascii.IsFont("standard") || !ascii.IsFont("shadow") || !ascii.IsFont("thinkertoy") {
		t.Error("expected 'standard', 'shadow', and 'thinkertoy' to be valid fonts")
	}
	if ascii.IsFont("invalid") {
		t.Error("expected 'invalid' not to be a valid font")
	}
}
