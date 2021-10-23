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
	rootDir := os.TempDir()
	indexpath := filepath.Join(rootDir, "index.db")
	db, err := openDB(indexpath)
	if err != nil {
		t.Fatalf("couldn't load Index: %s", err)
	}
	defer db.Close()
	if err := migrateDB(db); err != nil {
		t.Fatalf("couldn't migrate: %s", err)
	}
	content := "hello world!\n"
	r := strings.NewReader(content)
	id, err := createFile(db, rootDir, r)
	if err != nil {
		t.Fatalf("couldn't create file: %s", err)
	}

	f, err := getFile(db, rootDir, id)
	if err != nil {
		t.Fatal(err)
	}

	if f.ID != id {
		t.Fatalf("wrong id, exp: %d got %d", id, f.ID)
	}

	if f.Content != content {
		t.Fatalf("expected content :%s, got :%s", content, f.Content)
	}
	if f.ObjectID == "" {
		t.Fatal("expected objectID not to be empty")
	}
	if exp := "text/plain; charset=utf-8"; f.Type != exp {
		t.Fatalf("expected %s, got %s", exp, f.Type)
	}

	updatedContent := "hello new world!\n"
	ur := strings.NewReader(updatedContent)
	if err := updateFile(db, rootDir, id, ur); err != nil {
		t.Fatal(err)
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
