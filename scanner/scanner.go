package scanner

import (
	"fmt"
	"path/filepath"
)

const TopN = 20

var skipDirs = map[string]bool{
	"/proc": true, "/sys": true, "/dev": true,
	"/run":  true, "/tmp": true,
}

type FileEntry struct {
	Path string
	Size int64
}

type DupGroup struct {
	Size  int64
	Paths []string
}

func shouldSkip(path, root string) bool {
	if skipDirs[path] {
		return true
	}
	return path != root &&
		filepath.Dir(path) == root &&
		len(filepath.Base(path)) > 0 &&
		filepath.Base(path)[0] == '.'
}

func FormatSize(b int64) string {
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
