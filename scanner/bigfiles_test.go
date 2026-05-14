package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.WriteFile(path, make([]byte, size), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestFindLargestFiles_Order(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "small.txt"), 100)
	writeFile(t, filepath.Join(root, "big.txt"), 1024)
	writeFile(t, filepath.Join(root, "medium.txt"), 512)

	entries, err := FindLargestFiles(root, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	if entries[0].Size != 1024 {
		t.Errorf("want largest 1024, got %d", entries[0].Size)
	}
	if entries[1].Size != 512 {
		t.Errorf("want second 512, got %d", entries[1].Size)
	}
	if entries[2].Size != 100 {
		t.Errorf("want third 100, got %d", entries[2].Size)
	}
}

func TestFindLargestFiles_TopN(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), 300)
	writeFile(t, filepath.Join(root, "b.txt"), 200)
	writeFile(t, filepath.Join(root, "c.txt"), 100)

	entries, err := FindLargestFiles(root, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	if entries[0].Size != 300 || entries[1].Size != 200 {
		t.Errorf("unexpected top-2: %v", entries)
	}
}

func TestFindLargestFiles_SkipsHiddenTopLevelDir(t *testing.T) {
	root := t.TempDir()
	hidden := filepath.Join(root, ".hidden")
	if err := os.Mkdir(hidden, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(hidden, "secret.txt"), 9999)
	writeFile(t, filepath.Join(root, "visible.txt"), 1)

	entries, err := FindLargestFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if filepath.Dir(e.Path) == hidden {
			t.Errorf("hidden dir file should be skipped, got %q", e.Path)
		}
	}
}

func TestFindLargestFiles_SubdirFiles(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(sub, "deep.txt"), 2048)
	writeFile(t, filepath.Join(root, "shallow.txt"), 512)

	entries, err := FindLargestFiles(root, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	if entries[0].Size != 2048 {
		t.Errorf("want largest 2048, got %d", entries[0].Size)
	}
}
