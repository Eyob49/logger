# logger

A small Go logger that writes level-filtered messages to standard output, a file, or both.

## Install

Once this repository is public, install it with:

```sh
go get github.com/Eyob49/logger
```

The module path in `go.mod` matches this import path. The repository is currently private, so other users cannot install it until it is made public.

## Usage

```go
package main

import (
	"fmt"

	"github.com/Eyob49/logger"
)

func main() {
	l, err := logger.New(logger.ModeStdout, logger.LevelInfo, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	l.Info("server started")
}
```

For file output, provide a path:

```go
func main() {
	l, err := logger.New(logger.ModeFile, logger.LevelInfo, "app.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	l.Info("server started")
}
```

The path argument is ignored for `ModeStdout`.

Output looks like:

```text
2026-10-06 11:53:28 [INFO] server started
```

## Modes

- `ModeStdout` — write to standard output.
- `ModeFile` — write to the file at the provided path.
- `ModeBoth` — write to standard output and the file.

## Levels and methods

The logger provides `Debug`, `Info`, `Warn`, and `Error`. Messages below the chosen minimum level are dropped; for example, `LevelWarn` shows only `Warn` and `Error`.

Available levels: `LevelDebug`, `LevelInfo`, `LevelWarn`, and `LevelError`.

## Behavior

- Safe to use from multiple goroutines.
- Write errors are reported to standard error; they do not crash the program.
- `Close` releases the file, if one was opened. Repeated calls return `nil`.

See [NOTES.md](./NOTES.md) for design notes.
