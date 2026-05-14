package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type dupGroup struct {
	size  int64
	paths []string
}

type dupResultMsg struct {
	groups []dupGroup
	err    error
}

func scanDupFiles() tea.Cmd {
	return func() tea.Msg {
		groups, err := findDuplicateFiles("/")
		return dupResultMsg{groups: groups, err: err}
	}
}

func findDuplicateFiles(root string) ([]dupGroup, error) {
	skipDirs := map[string]bool{
		"/proc": true, "/sys": true, "/dev": true,
		"/run":  true, "/tmp": true,
	}

	// group paths by size first — cheap pre-filter
	bySize := make(map[int64][]string)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[path] {
				return filepath.SkipDir
			}
			if path != root && filepath.Dir(path) == root && len(filepath.Base(path)) > 0 && filepath.Base(path)[0] == '.' {
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

	// hash only files that share a size
	byHash := make(map[string][]string)
	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			h, herr := md5File(p)
			if herr != nil {
				continue
			}
			key := fmt.Sprintf("%d:%s", size, h)
			byHash[key] = append(byHash[key], p)
		}
	}

	var groups []dupGroup
	for key, paths := range byHash {
		if len(paths) < 2 {
			continue
		}
		var size int64
		fmt.Sscanf(key, "%d:", &size)
		sort.Strings(paths)
		groups = append(groups, dupGroup{size: size, paths: paths})
	}
	sort.Slice(groups, func(i, j int) bool {
		wi := groups[i].size * int64(len(groups[i].paths)-1)
		wj := groups[j].size * int64(len(groups[j].paths)-1)
		return wi > wj
	})
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
