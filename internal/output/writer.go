package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/kzelealem/apex/internal/types"
)

func ProcessAndWriteFiles(cfg *types.Config, paths []string, tree string) error {
	switch cfg.OutputFormat {
	case "markdown":
		return writeMarkdown(cfg, paths, tree)
	case "json":
		return writeJSON(cfg, paths, tree)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.OutputFormat)
	}
}

func writeMarkdown(cfg *types.Config, paths []string, tree string) error {
	outputFile, err := os.Create(cfg.OutputFileName)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFileName, err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	writer.WriteString("# Project Structure\n\n")
	writer.WriteString("```\n")
	writer.WriteString(tree)
	writer.WriteString("```\n\n")

	if !cfg.TreeOnly {
		writer.WriteString("# File Contents\n\n")
		return processFilesForMarkdown(cfg, writer, paths)
	}
	return nil
}

func writeJSON(cfg *types.Config, paths []string, tree string) error {
	output := make(map[string]interface{})
	output["projectTree"] = tree
	if !cfg.TreeOnly {
		fileContents := make(map[string]string)
		for _, path := range paths {
			relPath, _ := filepath.Rel(cfg.RootDir, path)
			content, err := os.ReadFile(path)
			if err != nil {
				fileContents[relPath] = fmt.Sprintf("Error reading file: %v", err)
			} else {
				fileContents[relPath] = string(content)
			}
		}
		output["fileContents"] = fileContents
	}
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.OutputFileName, jsonData, 0644)
}

func processFilesForMarkdown(cfg *types.Config, writer *bufio.Writer, paths []string) error {
	numWorkers := runtime.NumCPU()
	pathChan := make(chan string, len(paths))
	resultChan := make(chan types.FileData, len(paths))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go fileReaderWorker(&wg, pathChan, resultChan, cfg.RootDir)
	}

	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)
	wg.Wait()
	close(resultChan)

	for data := range resultChan {
		if !cfg.Quiet {
			fmt.Printf("  [read] %s\n", data.Path)
		}
		fmt.Fprintf(writer, "---\n\n### `%s`\n\n", data.Path)
		if data.Err != nil {
			fmt.Fprintf(writer, "```\nError reading file: %v\n```\n\n", data.Err)
			continue
		}
		lang := strings.TrimPrefix(filepath.Ext(data.Path), ".")
		fmt.Fprintf(writer, "```%s\n", lang)
		writer.Write(data.Content)
		writer.WriteString("\n```\n\n")
	}

	return nil
}

func fileReaderWorker(wg *sync.WaitGroup, pathChan <-chan string, resultChan chan<- types.FileData, rootDir string) {
	defer wg.Done()
	for path := range pathChan {
		relPath, _ := filepath.Rel(rootDir, path)
		content, err := os.ReadFile(path)
		resultChan <- types.FileData{
			Path:    relPath,
			Content: content,
			Err:     err,
		}
	}
}
