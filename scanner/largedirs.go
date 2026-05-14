package scanner

import (
	"os"
	"path/filepath"
	"sort"
)

func FindLargestDirs(root string, n int) ([]FileEntry, error) {
	sizes := make(map[string]int64)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkip(path, root) {
				return filepath.SkipDir
			}
			if _, exists := sizes[path]; !exists {
				sizes[path] = 0
			}
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		size := info.Size()
		dir := filepath.Dir(path)
		for {
			sizes[dir] += size
			if dir == root {
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
		return nil
	})

	entries := make([]FileEntry, 0, len(sizes))
	for path, size := range sizes {
		if path == root {
			continue
		}
		entries = append(entries, FileEntry{Path: path, Size: size})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Size > entries[j].Size })
	if len(entries) > n {
		entries = entries[:n]
	}
	return entries, err
}
