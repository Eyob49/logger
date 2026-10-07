package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Eyob49/logger"
)

func main() {
	dir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		fmt.Println("could not create temp dir:", err)
		return
	}
	defer os.RemoveAll(dir)
	testStdout()
	testFile(dir)
	testBoth(dir)
	testBadPath()
	testConcurrent(dir)
}

func showFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("could not read file:", err)
		return
	}
	fmt.Printf("file %s has %d lines\n", filepath.Base(path), strings.Count(string(data), "\n"))
}

func testStdout() {
	fmt.Println("\n--- stdout: Close, then log ---")
	l, err := logger.New(logger.ModeStdout, logger.LevelDebug, "x")
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	l.Info("before close")
	fmt.Println("Close returned:", l.Close())
	l.Info("after close")
}

func testFile(dir string) {
	fmt.Println("\n--- file: filtering and append ---")
	path := filepath.Join(dir, "test_file.log")
	l, err := logger.New(logger.ModeFile, logger.LevelWarn, path)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	l.Info("filtered out")
	l.Warn("kept")
	l.Error("kept")
	fmt.Println("Close returned:", l.Close())
	showFile(path)
}

func testBoth(dir string) {
	fmt.Println("\n--- both: Close, then log ---")
	path := filepath.Join(dir, "test_both.log")
	l, err := logger.New(logger.ModeBoth, logger.LevelDebug, path)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	l.Info("before close")
	fmt.Println("Close returned:", l.Close())
	l.Info("after close")
	showFile(path)
}

func testBadPath() {
	fmt.Println("\n--- bad path ---")
	l, err := logger.New(logger.ModeFile, logger.LevelInfo, "no_such_folder/app.log")
	fmt.Println("logger is nil:", l == nil, "| error:", err)
}

func testConcurrent(dir string) {
	fmt.Println("\n--- concurrent: 100 goroutines x 10 lines ---")
	path := filepath.Join(dir, "test_concurrent.log")
	l, err := logger.New(logger.ModeFile, logger.LevelInfo, path)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				l.Info(fmt.Sprintf("goroutine %d line %d", id, j))
			}
		}(i)
	}
	wg.Wait()
	l.Close()
	showFile(path)
}
