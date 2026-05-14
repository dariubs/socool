package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicateFiles_FindsDups(t *testing.T) {
	root := t.TempDir()
	content := []byte("duplicate content")
	os.WriteFile(filepath.Join(root, "a.txt"), content, 0644)
	os.WriteFile(filepath.Join(root, "b.txt"), content, 0644)
	os.WriteFile(filepath.Join(root, "c.txt"), []byte("unique"), 0644)

	groups, err := FindDuplicateFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("want 1 group, got %d", len(groups))
	}
	if len(groups[0].Paths) != 2 {
		t.Errorf("want 2 paths in group, got %d", len(groups[0].Paths))
	}
}

func TestFindDuplicateFiles_NoDups(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(root, "b.txt"), []byte("world"), 0644)

	groups, err := FindDuplicateFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Errorf("want 0 groups, got %d", len(groups))
	}
}

func TestFindDuplicateFiles_SkipsZeroByteFiles(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "empty1.txt"), []byte{}, 0644)
	os.WriteFile(filepath.Join(root, "empty2.txt"), []byte{}, 0644)

	groups, err := FindDuplicateFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Errorf("zero-byte files should not be grouped, got %d groups", len(groups))
	}
}

func TestFindDuplicateFiles_SameSizeDifferentContent(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.txt"), []byte("aaaa"), 0644)
	os.WriteFile(filepath.Join(root, "b.txt"), []byte("bbbb"), 0644)

	groups, err := FindDuplicateFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Errorf("same-size different-content files should not be grouped, got %d", len(groups))
	}
}

func TestFindDuplicateFiles_TopN(t *testing.T) {
	root := t.TempDir()
	for i := range 5 {
		// Each pair is a distinct duplicate group with distinct content
		content := make([]byte, i+1)
		for j := range content {
			content[j] = byte(i + 1)
		}
		os.WriteFile(filepath.Join(root, filepath.Join(root, "")+string(rune('a'+i*2))+".txt"), content, 0644)
		os.WriteFile(filepath.Join(root, filepath.Join(root, "")+string(rune('a'+i*2+1))+".txt"), content, 0644)
	}

	groups, err := FindDuplicateFiles(root, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) > 3 {
		t.Errorf("want at most 3 groups, got %d", len(groups))
	}
}

func TestFindDuplicateFiles_MultipleGroupsSortedByWaste(t *testing.T) {
	root := t.TempDir()

	// small group: 2 copies of 10 bytes = 10 bytes wasted
	small := []byte("0123456789")
	os.WriteFile(filepath.Join(root, "s1.txt"), small, 0644)
	os.WriteFile(filepath.Join(root, "s2.txt"), small, 0644)

	// large group: 2 copies of 100 bytes = 100 bytes wasted
	large := make([]byte, 100)
	for i := range large {
		large[i] = 0xFF
	}
	os.WriteFile(filepath.Join(root, "l1.txt"), large, 0644)
	os.WriteFile(filepath.Join(root, "l2.txt"), large, 0644)

	groups, err := FindDuplicateFiles(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("want 2 groups, got %d", len(groups))
	}
	if groups[0].Size != 100 {
		t.Errorf("largest-waste group should be first (size 100), got size %d", groups[0].Size)
	}
}
