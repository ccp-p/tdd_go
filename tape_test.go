package main

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))
// 123abc
	file.Seek(0, 0)
	newFileContents, _ := io.ReadAll(file)
	t.Log(string(newFileContents))
	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}