package types

type FileData struct {
	Path    string
	Content []byte
	Err     error
}

type Config struct {
	RootDir           string
	OutputFileName    string
	IgnoreFileName    string
	AdditionalIgnores []string
	IncludePatterns   []string
	TreeOnly          bool
	OutputFormat      string
	MaxFileSize       int64
	Quiet             bool
}
