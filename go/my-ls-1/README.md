## my-ls

A clean, Unix-style reimplementation of the ls command written in Go, focusing on correctness, clarity, and faithful behavior rather than shortcuts.

This project recreates core ls functionality step by step, closely matching real ls semantics while keeping the codebase understandable and modular.

### Features

- List files and directories

- Long listing format (`-l`)

- Show hidden files (`-a`)

- Recursive directory listing (`-R`)

- List directories themselves, not their contents (`-d`)

- One entry per line (`-1`)

- Reverse sorting order (`-r`)

- Sort by modification time, newest first (`-t`)

- Column output with terminal-width awareness

- Correct handling of `.` and `..`

- Accurate name sorting matching `ls` behavior

- Symlink display (`name -> target`)

- Directory block totals (`total N`)

- Realistic timestamp formatting (6-month rule)

- POSIX-style error handling

- Buffered output for improved performance

- UID/GID caching for efficient long listings

- Help output (`--help`)


### Usage
```
my-ls [OPTION]... [FILE]...
```

If no files are provided, the current directory (.) is listed.

### Options
Flag	Description

-l	    Use a long listing format

-a	    Do not ignore entries starting with .

-R	    List subdirectories recursively

-d	    List directories themselves, not their contents

-r      Reverse order while sorting

-t      Sort by modification time, newest first

-1	    List one file per line

--help	Display help and exit

### Behavior Notes

- Errors accessing files or directories are reported, but do not stop execution.

- Invalid flags are treated as fatal and exit immediately.

- Exit status matches ls semantics (non-zero if errors occurred).

- Output formatting is designed to closely resemble real ls, while keeping the implementation readable.

- Columns adapt to terminal width; narrow terminals fall back to single-column output.

- Absolute paths are only printed if provided by the user, matching ls behavior.

### Design

The project is split into clear layers:

- cmd — argument parsing and program entry

- internal — filesystem traversal and control flow

- commands — output formatting and presentation

- logger — error handling and help output

- models — shared data structures

This separation keeps logic, formatting, and error handling independent and easy to reason about.

### Philosophy

This is not a minimal clone or a wrapper around system calls.
The goal is to understand and reimplement how ls actually works:

No shelling writer.Out to ls

No shortcuts around filesystem semantics

Clear, explicit handling of edge cases

Step-by-step construction with correctness first

Build & Run
```
go run main.go
```

Or build a binary:
```
go build -o my-ls
./my-ls
```

### Contributors
- tdiridis