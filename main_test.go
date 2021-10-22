package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateDB(t *testing.T) {
	indexpath := filepath.Join(os.TempDir(), "index.db")
	db, err := openDB(indexpath)
	if err != nil {
		t.Fatalf("couldn't load Index: %s", err)
	}
	defer db.Close()
	if err := migrateDB(db); err != nil {
		t.Fatalf("couldn't migrate: %s", err)
	}
}

func TestCreateFile(t *testing.T) {
	indexpath := filepath.Join(os.TempDir(), "index.db")
	db, err := openDB(indexpath)
	if err != nil {
		t.Fatalf("couldn't load Index: %s", err)
	}
	defer db.Close()
	if err := migrateDB(db); err != nil {
		t.Fatalf("couldn't migrate: %s", err)
	}
	r := strings.NewReader("hello world!")
	tmp := os.TempDir()
	if _, err := createFile(db, tmp, r); err != nil {
		t.Fatalf("couldn't create file: %s", err)
	}
}

func TestGetObjectPath(t *testing.T) {
	tests := map[string]string{
		"aa2f368177a48ff6b1b8304d21ca584629c57c8a": "root/objects/aa/2f368177a48ff6b1b8304d21ca584629c57c8a",
		"da39a3ee5e6b4b0d3255bfef95601890afd80709": "root/objects/da/39a3ee5e6b4b0d3255bfef95601890afd80709",
	}

	for test, exp := range tests {
		got := getObjectPath("root", test)
		if got != exp {
			t.Errorf("expected: %s, got %s", exp, got)
		}
	}
}
