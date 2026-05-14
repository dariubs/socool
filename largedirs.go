package main

import (
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type dirsResultMsg struct {
	dirs []fileEntry
	err  error
}

func scanLargestDirs() tea.Cmd {
	return func() tea.Msg {
		dirs, err := findLargestDirs("/", topN)
		return dirsResultMsg{dirs: dirs, err: err}
	}
}

func findLargestDirs(root string, n int) ([]fileEntry, error) {
	sizes := make(map[string]int64)

	skipDirs := map[string]bool{
		"/proc": true, "/sys": true, "/dev": true,
		"/run":  true, "/tmp": true,
	}

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

	entries := make([]fileEntry, 0, len(sizes))
	for path, size := range sizes {
		if path == root {
			continue
		}
		entries = append(entries, fileEntry{path: path, size: size})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].size > entries[j].size })
	if len(entries) > n {
		entries = entries[:n]
	}
	return entries, err
}
