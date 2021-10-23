package main

import (
	"os"
	"strings"
	"testing"
)

func TestStoreMethods(t *testing.T) {
	s, err := NewStore(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := migrateDB(s.DB); err != nil {
		t.Fatalf("couldn't migrate: %s", err)
	}
	content := "hello world!\n"
	r := strings.NewReader(content)
	id, err := s.Create(r)
	if err != nil {
		t.Fatalf("couldn't create file: %s", err)
	}

	f, err := s.Get(id)
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
	if err := s.Update(id, ur); err != nil {
		t.Fatal(err)
	}

	f1, err := s.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if f1.Content != updatedContent {
		t.Fatalf("expected content :%s, got :%s", updatedContent, f1.Content)
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
