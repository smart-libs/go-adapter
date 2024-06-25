package test

import (
	"bytes"
	"io"
	"os"
)

func CaptureStdout(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	_ = w.Close()
	os.Stdout = old // restoring the real stdout
	return <-outC
}

func CaptureStderr(f func()) string {
	old := os.Stderr // keep backup of the real stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stderr = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	_ = w.Close()
	os.Stderr = old // restoring the real stdout
	return <-outC
}
