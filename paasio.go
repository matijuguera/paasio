package paasio

import (
	"io"
	"sync"
)

type MyReadCounter struct {
	io.Reader
	Counter
}

type MyWriteCounter struct {
	io.Writer
	Counter
}

type Counter struct {
	bytes      int64
	operations int
	sync.RWMutex
}

type MyReadWriteCounter struct {
	WriteCounter
	ReadCounter
}

func NewReadCounter(r io.Reader) ReadCounter {
	return &MyReadCounter{r, Counter{}}
}

func NewWriteCounter(w io.Writer) WriteCounter {
	return &MyWriteCounter{w, Counter{}}
}

func NewReadWriteCounter(rw io.ReadWriter) ReadWriteCounter {
	return MyReadWriteCounter{
		NewWriteCounter(rw), NewReadCounter(rw),
	}
}

func (r *MyReadCounter) ReadCount() (n int64, nops int) {
	return r.Count()
}

func (rc *MyReadCounter) Read(b []byte) (int, error) {
	m, err := rc.Reader.Read(b)
	rc.AddBytes(m)
	return m, err
}

func (w *MyWriteCounter) WriteCount() (n int64, nops int) {
	return w.Count()
}

func (wc *MyWriteCounter) Write(b []byte) (int, error) {
	m, err := wc.Writer.Write(b)
	wc.AddBytes(m)
	return m, err
}

func (c *Counter) Count() (int64, int) {
	c.RLock()
	defer c.RUnlock()
	return c.bytes, c.operations
}

func (c *Counter) AddBytes(n int) {
	c.Lock()
	defer c.Unlock()
	c.bytes += int64(n)
	c.operations++
}
