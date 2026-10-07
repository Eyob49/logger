package logger

import (
	"path/filepath"
	"testing"
)

func TestCloseIsIdempotent(t *testing.T) {
	l, err := New(ModeFile, LevelInfo, filepath.Join(t.TempDir(), "test.log"))
	if err != nil {
		t.Fatal(err)
	}

	if err := l.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}
