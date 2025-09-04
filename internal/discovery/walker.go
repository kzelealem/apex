package discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/kzelealem/apex/internal/types"
	gitignore "github.com/sabhiram/go-gitignore"
)

func DiscoverFilesAndBuildTree(cfg *types.Config, ignoreMatcher *gitignore.GitIgnore) ([]string, string, error) {
	var filePaths []string
	var treeBuilder strings.Builder

	walkErr := filepath.WalkDir(cfg.RootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cfg.RootDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			base := filepath.Base(cfg.RootDir)
			treeBuilder.WriteString(fmt.Sprintf("%s\n", base))
			if !cfg.Quiet {
				fmt.Printf("[scan] dir: %s\n", base)
			}
			return nil
		}

		matchPath := relPath
		if d.IsDir() {
			matchPath += string(filepath.Separator)
		}
		if ignoreMatcher.MatchesPath(matchPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if len(cfg.IncludePatterns) > 0 {
			included := false
			for _, pattern := range cfg.IncludePatterns {
				if matched, _ := filepath.Match(pattern, d.Name()); matched {
					included = true
					break
				}
			}
			if !d.IsDir() && !included {
				return nil
			}
		}

		if !d.IsDir() && cfg.MaxFileSize > 0 {
			info, err := d.Info()
			if err == nil && info.Size() > cfg.MaxFileSize {
				if !cfg.Quiet {
					fmt.Printf("  [scan] skip: %s (size limit exceeded)\n", relPath)
				}
				return nil
			}
		}

		depth := strings.Count(relPath, string(filepath.Separator))
		indent := strings.Repeat("│   ", depth)
		prefix := "├── "
		treeBuilder.WriteString(fmt.Sprintf("%s%s%s\n", indent, prefix, d.Name()))

		if !d.IsDir() {
			filePaths = append(filePaths, path)
			if !cfg.Quiet {
				fmt.Printf("  [scan] file: %s\n", relPath)
			}
		} else {
			if !cfg.Quiet {
				fmt.Printf("  [scan] dir: %s\n", relPath)
			}
		}
		return nil
	})

	if walkErr == filepath.SkipDir {
		return filePaths, treeBuilder.String(), nil
	}

	return filePaths, treeBuilder.String(), walkErr
}
