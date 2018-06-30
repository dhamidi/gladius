# Description

![[Documentation](https://godoc.org/github.com/dhamidi/gladius?status.svg)]

Gladius is a headless, scriptable text editor for use in other programs.

*Use cases*:

- Automate complex editing of files based on a well known mental model
- Provide a backend for text editing operations in other programs (e.g. a line editor)

# Status

- 2018-06-30: Basic editing operations are implemented.

# Roadmap

- Add more movement commands to efficiently navigate a buffer
- Add regexp search (forward and backward)
- Copy / Cut / Paste text into registers
- Add line-based command language for simple batch editing
- Add file management to make common open-edit-save cycle easy

# Usage

Currently only usage in Go programs is supported (see example):

```go
package main

import (
        "fmt"
        "github.com/dhamidi/gladius"
)

func main() {
        editor := gladius.NewEditor("some initial text")
        editor.
                Delete(len("some initial ")).
                Insert("edited ")

        fmt.Printf("Editor contains: %s\n", editor.String())
        // Editor contains: edited text
}

```

## Creating a new editor

An editor operates on a text buffer (based on a piece table) and maintains a position in that buffer.

You can create a new editor instance by calling `gladius.NewEditor("initial text")`.

The cursor position of a new editor instance is always at the beginning of the text (`0`).


## API overview

| Operation                    | Effect                                                    |
| ---------                    | ------                                                    |
| `editor.Insert("text")`      | inserts `"text"` at the cursor position                   |
| `editor.Delete(3)`           | deletes `3` characters right of the cursor                |
| `editor.String()`            | returns the contents of the underlying buffer as a string |
| `editor.Forward(1)`          | moves cursor forward by `1` character                     |
| `editor.Backward(1)`         | moves cursor backward by `1` character                    |
| `editor.BeginningOfBuffer()` | Moves cursor to beginning of text                         |
| `editor.EndOfBuffer()`       | Moves cursor to end of text                               |
| `editor.BeginningOfLine()`   | Moves cursor to beginning of the current line             |
| `editor.EndOfLine()`         | Moves cursor to end of the current line                   |
