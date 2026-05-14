package scanner

import (
	"os"
	"path/filepath"
	"sort"
)

func FindLargestFiles(root string, n int) ([]FileEntry, error) {
	var top []FileEntry
	minSize := int64(0)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkip(path, root) {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		size := info.Size()
		if len(top) < n || size > minSize {
			top = append(top, FileEntry{Path: path, Size: size})
			sort.Slice(top, func(i, j int) bool { return top[i].Size > top[j].Size })
			if len(top) > n {
				top = top[:n]
			}
			if len(top) == n {
				minSize = top[n-1].Size
			}
		}
		return nil
	})

	return top, err
}
