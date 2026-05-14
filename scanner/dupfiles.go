package scanner

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type hashKey struct {
	size int64
	hash string
}

func FindDuplicateFiles(root string, n int) ([]DupGroup, error) {
	bySize := make(map[int64][]string)

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
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			return nil
		}
		bySize[info.Size()] = append(bySize[info.Size()], path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	byHash := make(map[hashKey][]string)
	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			h, herr := md5File(p)
			if herr != nil {
				continue
			}
			key := hashKey{size: size, hash: h}
			byHash[key] = append(byHash[key], p)
		}
	}

	var groups []DupGroup
	for key, paths := range byHash {
		if len(paths) < 2 {
			continue
		}
		sort.Strings(paths)
		groups = append(groups, DupGroup{Size: key.size, Paths: paths})
	}
	sort.Slice(groups, func(i, j int) bool {
		wi := groups[i].Size * int64(len(groups[i].Paths)-1)
		wj := groups[j].Size * int64(len(groups[j].Paths)-1)
		return wi > wj
	})
	if len(groups) > n {
		groups = groups[:n]
	}
	return groups, nil
}

func md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
