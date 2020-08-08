package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

//Using fake Stdin/Stdout for testing.
func main() {
	r, w, err := os.Pipe()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Unable to create a new pipe")
		os.Exit(2)
	}

	var (
		orgStdout *os.File
		orgStdin  *os.File
	)
	orgStdout = os.Stdout
	orgStdin = os.Stdin

	defer func() {
		os.Stdout = orgStdout
		os.Stdin = orgStdin
	}()

	os.Stdout = w
	os.Stdin = r
	_, err = io.Copy(os.Stdout, bytes.NewBuffer([]byte("how are you? \ni am fine, thanks for your care")))
	w.Close()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "copy failed")
		os.Exit(2)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		_, _ = fmt.Fprintln(os.Stderr, scanner.Text())
	}

}
