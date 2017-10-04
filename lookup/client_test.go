package lookup

import (
	"testing"
)

func TestExists(t *testing.T) {
	client := NewMemoryStorage()
	client.Add("ABCDE")
	exists := client.Exists("ABCDE")
	if !exists {
		t.Error("ABCDE should exist but was not found")
	}
}

func TestNotExists(t *testing.T) {
	client := NewMemoryStorage()
	exists := client.Exists("ABCDE")
	if exists {
		t.Error("ABCDE should not exist but was found")
	}
}
