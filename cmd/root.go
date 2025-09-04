package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kzelealem/apex/internal/discovery"
	"github.com/kzelealem/apex/internal/output"
	"github.com/kzelealem/apex/internal/types"
	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

var cfg types.Config

var rootCmd = &cobra.Command{
	Use:   "apex [directory]",
	Short: "A fast and flexible tool for generating project structure documentation.",
	Long: `Apex performs an efficient scan of a directory to generate documentation
in various formats. It supports advanced filtering and customization to provide
a comprehensive project overview.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cfg.RootDir = args[0]
		} else {
			cfg.RootDir = "."
		}
		if err := run(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.OutputFileName, "output", "o", "project_structure", "Name of the output file.")
	rootCmd.Flags().StringVarP(&cfg.IgnoreFileName, "ignore-file", "i", ".gitignore", "Path to a custom ignore file.")
	rootCmd.Flags().StringSliceVarP(&cfg.AdditionalIgnores, "add-ignore", "a", []string{}, "Additional patterns to ignore.")
	rootCmd.Flags().StringSliceVar(&cfg.IncludePatterns, "include", []string{}, "Only include files matching these patterns (e.g., '*.go').")
	rootCmd.Flags().BoolVarP(&cfg.TreeOnly, "tree-only", "t", false, "Only output the directory tree structure.")
	rootCmd.Flags().StringVarP(&cfg.OutputFormat, "format", "f", "markdown", "Output format (markdown, json).")
	rootCmd.Flags().Int64Var(&cfg.MaxFileSize, "max-size", 0, "Maximum file size to include (in bytes).")
	rootCmd.Flags().BoolVarP(&cfg.Quiet, "quiet", "q", false, "Suppress all interactive output except for errors.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("output") {
		switch cfg.OutputFormat {
		case "json":
			cfg.OutputFileName += ".json"
		case "markdown":
			fallthrough
		default:
			cfg.OutputFileName += ".md"
		}
	}

	if !cfg.Quiet {
		fmt.Printf("[+] Starting Apex scan of: %s\n\n", cfg.RootDir)
	}

	cfg.AdditionalIgnores = append(cfg.AdditionalIgnores, cfg.OutputFileName, ".git")

	ignoreMatcher, err := loadIgnoreMatcher(cfg.RootDir, cfg.IgnoreFileName, cfg.AdditionalIgnores)
	if err != nil {
		return fmt.Errorf("failed to load ignore rules: %w", err)
	}

	if !cfg.Quiet {
		fmt.Println("[1/3] Scanning project structure...")
	}
	filePaths, tree, err := discovery.DiscoverFilesAndBuildTree(&cfg, ignoreMatcher)
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if !cfg.TreeOnly && !cfg.Quiet {
		fmt.Println("\n[2/3] Reading file contents...")
	}

	if !cfg.Quiet {
		step := "3/3"
		if cfg.TreeOnly {
			step = "2/2"
		}
		fmt.Printf("\n[%s] Writing output to %s...\n", step, cfg.OutputFileName)
	}

	if err := output.ProcessAndWriteFiles(&cfg, filePaths, tree); err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	if !cfg.Quiet {
		fmt.Printf("\n[+] Success! Project documentation saved to '%s' in %s format.\n", cfg.OutputFileName, cfg.OutputFormat)
	}

	return nil
}

func loadIgnoreMatcher(rootDir, ignoreFileName string, additionalIgnores []string) (*gitignore.GitIgnore, error) {
	ignoreFilePath := filepath.Join(rootDir, ignoreFileName)
	ignoreMatcher, err := gitignore.CompileIgnoreFileAndLines(ignoreFilePath, additionalIgnores...)
	if err != nil {
		if os.IsNotExist(err) {
			if !cfg.Quiet {
				fmt.Printf("[!] Note: %s not found. Proceeding without it.\n", ignoreFileName)
			}
			return gitignore.CompileIgnoreLines(additionalIgnores...), nil
		}
		return nil, err
	}
	return ignoreMatcher, nil
}
