# ASCII Art Terminal Renderer

This Go program renders ASCII art for strings in various font styles, with customizable alignment options. It allows users to control the alignment of the rendered ASCII art with respect to the terminal window using the `--align` flag. The program supports four alignment types: **center**, **left**, **right**, and **justify**.

## Features
- **Align text** in four different ways: `center`, `left`, `right`, and `justify`.
- Supports multiple **font styles**, such as `standard`, `shadow`, etc.
- Automatically adjusts the graphical representation to fit the current terminal size.
- Handles single-string inputs as well as multi-line strings.
- Can output directly to the terminal or to a specified text file.
  
## Requirements
- Go 1.16 or higher.
  
## Installation

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/yourusername/ascii-art-terminal.git
    cd ascii-art-terminal
    ```

2. Build the project:

    ```bash
    go build .
    ```

## Usage

You can run the program with various flags and arguments.

### Basic Command Format

```bash
go run . [OPTION] [STRING] [BANNER]
```

## Flags
- --align: Specify the alignment type (center, left, right, justify).

## Examples
Center alignment:

```bash
go run . --align=center "hello" standard
```
This will print "hello" using the standard font, centered on the screen.

Output (example):

```
                                _                _    _
                               | |              | |  | |
                               | |__      ___   | |  | |    ___
                               |  _ \    / _ \  | |  | |   / _ \
                               | | | |  |  __/  | |  | |  | (_) |
                               |_| |_|   \___|  |_|  |_|   \___/
```
Left alignment:

```bash
go run . --align=left "Hello There" standard
```
This will print "Hello There" using the standard font, aligned to the left of the terminal.

Output (example):
``` _    _           _    _                 _______   _
|| |  | |         | |  | |               |__   __| | |
|| |__| |   ___   | |  | |    ___           | |    | |__      ___    _ __     ___
||  __  |  / _ \  | |  | |   / _ \          | |    |  _ \    / _ \  | '__|   / _ \
|| |  | | |  __/  | |  | |  | (_) |         | |    | | | |  |  __/  | |     |  __/
||_|  |_|  \___|  |_|  |_|   \___/          |_|    |_| |_|   \___|  |_|      \___/
```
Right alignment:

```bash
go run . --align=right "hello" shadow
```
This will print "hello" using the shadow font, right-aligned on the terminal.

Output (example):
```
                                                                                                                      
                                                                                      _|                _| _|          
                                                                                      _|_|_|     _|_|   _| _|   _|_|   
                                                                                      _|    _| _|_|_|_| _| _| _|    _| |
                                                                                      _|    _| _|       _| _| _|    _| |
                                                                                      _|    _|   _|_|_| _| _|   _|_|   
```
Justify alignment:

```bash
go run . --align=justify "how are you" shadow
```
This will print "how are you" using the shadow font, with words justified to fill the terminal width.

```
_|                                                                                                                         
_|_|_|     _|_|   _|      _|      _|                  _|_|_| _|  _|_|   _|_|                    _|    _|   _|_|   _|    _| |
_|    _| _|    _| _|      _|      _|                _|    _| _|_|     _|_|_|_|                  _|    _| _|    _| _|    _| |
_|    _| _|    _|   _|  _|  _|  _|                  _|    _| _|       _|                        _|    _| _|    _| _|    _| |
_|    _|   _|_|       _|      _|                      _|_|_| _|         _|_|_|                    _|_|_|   _|_|     _|_|_| |
```
## Optional Flags
- --output=<file>: Write the ASCII art output to a text file instead of printing it to the terminal.

Example:

```bash
go run . --align=center "hello" standard --output=output.txt
```
This will create a file output.txt containing the ASCII art output.

## Usage Message
If the alignment flag is not used correctly, the program will show the following usage message:
```
Usage: go run . [OPTION] [STRING] [BANNER]

Example: go run . --align=right something standard
```

## Notes
- The program dynamically adapts to the size of your terminal window, and will adjust the width of the output based on the terminal’s current width.
- The font style ([BANNER]) can be one of the predefined font styles, such as standard, shadow, etc.
- If no alignment flag (--align) is provided, the default alignment is left-aligned.
## Testing
You can run unit tests to validate the functionality of the program:

```bash
 go test ./test_files/ascii_test.go
 go test ./test_files/args_test.go
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.

---





