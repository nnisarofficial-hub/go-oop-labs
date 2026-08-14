package main

import (
	"fmt"
	"io"
	"strings"
)

type UpperCaseReader struct {
	source io.Reader
}

func NewUpperCaseReader(r io.Reader) *UpperCaseReader {
	return &UpperCaseReader{source: r}
}

func (u *UpperCaseReader) Read(p []byte) (int, error) {
	n, err := u.source.Read(p)
	for i := 0; i < n; i++ {
		if p[i] >= 'a' && p[i] <= 'z' {
			p[i] -= 32
		}
	}
	return n, err
}

type LineCountReader struct {
	source io.Reader
	Lines  int
}

func NewLineCountReader(r io.Reader) *LineCountReader {
	return &LineCountReader{
		source: r,
		Lines:  0,
	}
}
func (l *LineCountReader) Read(p []byte) (int, error) {
	n, err := l.source.Read(p)
	for i := 0; i < n; i++ {
		if p[i] == '\n' {
			l.Lines++
		}
	}
	return n, err
}

func main() {
	text := "hello world\ngo is great\ninterfaces are powerful\n"
	source := strings.NewReader(text)
	lineCounter := NewLineCountReader(source)
	upper := NewUpperCaseReader(lineCounter)
	result, _ := io.ReadAll(upper)
	fmt.Println(string(result))
	fmt.Printf("Line count: %d\n", lineCounter.Lines)
}
