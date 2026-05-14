package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindLargestDirs_Order(t *testing.T) {
	root := t.TempDir()

	a := filepath.Join(root, "a")
	b := filepath.Join(root, "b")
	for _, d := range []string{a, b} {
		if err := os.Mkdir(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(a, "big.txt"), 1024)
	writeFile(t, filepath.Join(b, "small.txt"), 100)

	entries, err := FindLargestDirs(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("want entries, got none")
	}
	if entries[0].Path != a {
		t.Errorf("want largest dir %q, got %q", a, entries[0].Path)
	}
	if entries[0].Size != 1024 {
		t.Errorf("want size 1024, got %d", entries[0].Size)
	}
}

func TestFindLargestDirs_ExcludesRoot(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(sub, "f.txt"), 512)

	entries, err := FindLargestDirs(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Path == root {
			t.Error("root itself should not appear in results")
		}
	}
}

func TestFindLargestDirs_CumulativeSizes(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "parent")
	child := filepath.Join(parent, "child")
	for _, d := range []string{parent, child} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(child, "f.txt"), 500)
	writeFile(t, filepath.Join(parent, "g.txt"), 300)

	entries, err := FindLargestDirs(root, 10)
	if err != nil {
		t.Fatal(err)
	}

	sizeOf := func(path string) int64 {
		for _, e := range entries {
			if e.Path == path {
				return e.Size
			}
		}
		t.Errorf("path %q not found in results", path)
		return -1
	}

	// parent accumulates both its own file and the child's file
	if got := sizeOf(parent); got != 800 {
		t.Errorf("parent size: want 800, got %d", got)
	}
	if got := sizeOf(child); got != 500 {
		t.Errorf("child size: want 500, got %d", got)
	}
}

func TestFindLargestDirs_TopN(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"x", "y", "z"} {
		d := filepath.Join(root, name)
		os.Mkdir(d, 0755)
		writeFile(t, filepath.Join(d, "f.txt"), 100)
	}

	entries, err := FindLargestDirs(root, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("want 2 entries, got %d", len(entries))
	}
}
