package paasio

import (
	"io"
	"sync"
)

type MyWriteCounter struct {
	io.Writer
	counter
}

func (w *MyWriteCounter) WriteCount() (n int64, nops int) {
	return w.count()
}

func (rc *MyWriteCounter) Write(p []byte) (int, error) {
	m, err := rc.Writer.Write(p)
	rc.addBytes(m)
	return m, err
}

type MyReadCounter struct {
	io.Reader
	counter
}

func (r *MyReadCounter) ReadCount() (n int64, nops int) {
	return r.count()
}

func (rc *MyReadCounter) Read(p []byte) (int, error) {
	m, err := rc.Reader.Read(p)
	rc.addBytes(m)
	return m, err
}

type MyReadWriteCounter struct {
	WriteCounter
	ReadCounter
}

type counter struct {
	bytes int64
	ops   int
	sync.RWMutex
}

func (c *counter) count() (int64, int) {
	c.RLock()
	defer c.RUnlock()
	return c.bytes, c.ops
}

func (c *counter) addBytes(n int) {
	c.Lock()
	defer c.Unlock()
	c.bytes += int64(n)
	c.ops++
}

func NewWriteCounter(w io.Writer) WriteCounter {
	return &MyWriteCounter{w, counter{}}
}

func NewReadCounter(r io.Reader) ReadCounter {
	return &MyReadCounter{r, counter{}}
}

func NewReadWriteCounter(rw io.ReadWriter) ReadWriteCounter {
	return MyReadWriteCounter{
		NewWriteCounter(rw), NewReadCounter(rw),
	}
}
