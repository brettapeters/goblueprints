package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events through code.
type Tracer interface {
	Trace(...interface{})
}

// New creates an instance of Tracer from an io.Writer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// Off creates a tracer that does nothing
func Off() Tracer {
	return &nilTracer{}
}

type nilTracer struct{}

func (n *nilTracer) Trace(a ...interface{}) {}
