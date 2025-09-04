# Apex

A handy CLI tool to quickly document the structure of your code repositories.

`apex` zips through a directory and spits out a single file showing your project's file tree and the contents of each file. It's fast, super customizable, and great for creating project overviews, preparing code for LLMs, or just getting a high-level look at a new codebase.

## What's cool about it?

-   **Super Fast:** It scans and reads files concurrently, so it doesn't drag on large projects.
-   **Smart Filtering:** You can use a `.gitignore` file, a different ignore file (like `.dockerignore`), or just add ignore patterns on the fly. You can also *only* include certain files or skip ones that are too big.
-   **Multiple Formats:** Get your output in Markdown (`.md`) for reading or JSON (`.json`) if you need to use the data in another script.
-   **Choose Your Vibe:** Watch it work with a live, step-by-step log, or run it with the `--quiet` flag in your build scripts.

## Getting Apex

### With Go

If you have the Go toolchain set up, this is the easiest way:

```sh
go install github.com/kzelealem/apex@latest
```

### From the Source

You can also just clone the repo and build it yourself.

```sh
git clone https://github.com/kzelealem/apex.git
cd apex
go build
./apex --help
```

## How to Use It

### The Basics

Just run `apex` in a project directory. It will create a `project_structure.md` file right there.

```sh
# Scan the current directory
apex

# Scan a different directory
apex ../my-other-project
```

### A More Advanced Example

Let's say you want to document a Go project. You only care about the `.go` and `.mod` files, you want to name the output `go_src.md`, and you want to ignore any file bigger than 50KB.

```sh
apex . -o go_src.md --include "*.go" --include "*.mod" --max-size 51200
```

### Creating JSON Output

To get a JSON file instead of Markdown, just change the format. `apex` is smart enough to change the file extension for you.

```sh
# This will create 'project_structure.json'
apex --format json
```

## Example Output

Wondering what the output looks like? Check out the example files included in this repository:

-   [`example_run.md`](./example_run.md): A sample of the default Markdown output.
-   [`example_run.json`](./example_run.json): A sample of the JSON output.

## All The Options

| Flag            | Short | Description                                         | Default                                |
| --------------- | ----- | --------------------------------------------------- | -------------------------------------- |
| `--output`      | `-o`  | Name of the output file.                            | `project_structure` + extension        |
| `--ignore-file` | `-i`  | Path to a custom ignore file.                       | `.gitignore`                           |
| `--add-ignore`  | `-a`  | Add an ignore pattern on the fly.                   | `[]`                                   |
| `--include`     |       | Only include files matching these patterns.         | `[]`                                   |
| `--max-size`    |       | Max file size in bytes (0 = no limit).              | `0`                                    |
| `--format`      | `-f`  | Output format (`markdown` or `json`).               | `markdown`                             |
| `--tree-only`   | `-t`  | Only show the file tree, not the file contents.     | `false`                                |
| `--quiet`       | `-q`  | Hide all the live output.                           | `false`                                |
| `--help`        | `-h`  | Show this help message.                             |                                        |

## Contributing

Got ideas? Feel free to fork the repo and submit a pull request.
