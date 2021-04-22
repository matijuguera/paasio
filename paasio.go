package paasio

import "io"

type MyWriteCounter struct {
	io.Writer
}

func (w MyWriteCounter) WriteCount() (n int64, nops int) {
	return 0, 0
}

func NewWriteCounter(w io.Writer) WriteCounter {
	return MyWriteCounter{w}
}
