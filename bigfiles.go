package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

const topN = 20

type fileEntry struct {
	path string
	size int64
}

type filesResultMsg struct {
	files []fileEntry
	err   error
}

func scanBiggestFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := findLargestFiles("/", topN)
		return filesResultMsg{files: files, err: err}
	}
}

func findLargestFiles(root string, n int) ([]fileEntry, error) {
	var top []fileEntry
	minSize := int64(0)

	skipDirs := map[string]bool{
		"/proc": true, "/sys": true, "/dev": true,
		"/run":  true, "/tmp": true,
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable paths
		}

		if d.IsDir() {
			if skipDirs[path] {
				return filepath.SkipDir
			}
			// skip hidden dirs at root level (e.g. /.Spotlight-V100)
			if path != root && len(filepath.Base(path)) > 0 && filepath.Base(path)[0] == '.' {
				base := filepath.Base(path)
				_ = base
				// only skip top-level hidden dirs
				if filepath.Dir(path) == root {
					return filepath.SkipDir
				}
			}
			return nil
		}

		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}

		size := info.Size()
		if len(top) < n || size > minSize {
			top = append(top, fileEntry{path: path, size: size})
			sort.Slice(top, func(i, j int) bool { return top[i].size > top[j].size })
			if len(top) > n {
				top = top[:n]
			}
			if len(top) == n {
				minSize = top[n-1].size
			}
		}
		return nil
	})

	return top, err
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
