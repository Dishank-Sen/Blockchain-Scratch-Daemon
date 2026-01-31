package utils

import (
	"os"
	"testing"
)

func TestSaveAndGetID(t *testing.T) {
	// cleanup before & after
	path := "/tmp/blocd.json"
	_ = os.Remove(path)
	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	expectedID := "node-123"

	// save
	if err := SaveID(expectedID); err != nil {
		t.Fatalf("SaveID failed: %v", err)
	}

	// get
	id, err := GetID()
	if err != nil {
		t.Fatalf("GetID failed: %v", err)
	}

	if id != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, id)
	}
}

func TestGetID_FileDoesNotExist(t *testing.T) {
	_ = os.Remove("/tmp/blocd.json")

	_, err := GetID()
	if err == nil {
		t.Fatalf("expected error when file does not exist")
	}
}
