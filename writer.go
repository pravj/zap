// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zap

import (
	"io"
	"io/ioutil"
)

// Discard is a convenience wrapper around ioutil.Discard.
var Discard = AddSync(ioutil.Discard)

// A WriteFlusher is an io.Writer that can also flush any buffered data.
type WriteFlusher interface {
	io.Writer
	Flush() error
}

// A WriteSyncer is an io.Writer that can also flush any buffered data. Note
// that *os.File (and thus, os.Stderr and os.Stdout) implement WriteSyncer.
type WriteSyncer interface {
	io.Writer
	Sync() error
}

// AddSync converts an io.Writer to a WriteSyncer. It attempts to be
// intelligent: if the concrete type of the io.Writer implements WriteSyncer or
// WriteFlusher, we'll use the existing Sync or Flush methods. If it doesn't,
func AddSync(w io.Writer) WriteSyncer {
	switch w.(type) {
	case WriteSyncer:
		// The concrete type is already a WriteSyncer (e.g., an *os.File).
		return w.(WriteSyncer)
	case WriteFlusher:
		// The concrete type implements a suitable Flush method, which we'll
		// just use instead of Sync.
		return flusherWrapper{w.(WriteFlusher)}
	default:
		// Fall back to a no-op Sync.
		return writerWrapper{w}
	}
}

type writerWrapper struct {
	io.Writer
}

func (w writerWrapper) Sync() error {
	return nil
}

type flusherWrapper struct {
	WriteFlusher
}

func (f flusherWrapper) Sync() error {
	return f.Flush()
}