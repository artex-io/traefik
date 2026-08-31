package metrics

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/metrics"
)

// CounterWithHeaders represents a counter that can use http.Header values as label values.
type CounterWithHeaders interface {
	Add(delta float64)
	With(headers http.Header, labelValues ...string) CounterWithHeaders
}

// MultiCounterWithHeaders collects multiple individual CounterWithHeaders and treats them as a unit.
type MultiCounterWithHeaders []CounterWithHeaders

// NewMultiCounterWithHeaders returns a multi-counter, wrapping the passed CounterWithHeaders.
func NewMultiCounterWithHeaders(c ...CounterWithHeaders) MultiCounterWithHeaders {
	return c
}

// Add adds the given delta value to the counter value.
func (c MultiCounterWithHeaders) Add(delta float64) {
	for _, counter := range c {
		counter.Add(delta)
	}
}

// With creates a new counter by appending the given label values and http.Header as labels and returns it.
func (c MultiCounterWithHeaders) With(headers http.Header, labelValues ...string) CounterWithHeaders {
	next := make(MultiCounterWithHeaders, len(c))
	for i := range c {
		next[i] = c[i].With(headers, labelValues...)
	}
	return next
}

// NewCounterWithNoopHeaders returns a CounterWithNoopHeaders.
func NewCounterWithNoopHeaders(counter metrics.Counter) CounterWithNoopHeaders {
	return CounterWithNoopHeaders{counter: counter}
}

// CounterWithNoopHeaders is a counter that satisfies CounterWithHeaders but ignores the given http.Header.
type CounterWithNoopHeaders struct {
	counter metrics.Counter
}

// Add adds the given delta value to the counter value.
func (c CounterWithNoopHeaders) Add(delta float64) {
	c.counter.Add(delta)
}

// With creates a new counter by appending the given label values and returns it.
func (c CounterWithNoopHeaders) With(_ http.Header, labelValues ...string) CounterWithHeaders {
	return NewCounterWithNoopHeaders(c.counter.With(labelValues...))
}

// NewMultiHistogramWithHeaders returns a multi-histogram, wrapping the passed ScalableHistogramWithHeaders.
func NewMultiHistogramWithHeaders(h ...ScalableHistogramWithHeaders) MultiHistogramWithHeaders {
	return h
}

// MultiHistogramWithHeaders collects multiple individual histograms and treats them as a unit.
type MultiHistogramWithHeaders []ScalableHistogramWithHeaders

// ObserveFromStart implements ScalableHistogramWithHeaders.
func (h MultiHistogramWithHeaders) ObserveFromStart(start time.Time) {
	for _, histogram := range h {
		histogram.ObserveFromStart(start)
	}
}

// Observe implements ScalableHistogramWithHeaders.
func (h MultiHistogramWithHeaders) Observe(v float64) {
	for _, histogram := range h {
		histogram.Observe(v)
	}
}

// With implements ScalableHistogramWithHeaders.
func (h MultiHistogramWithHeaders) With(headers http.Header, labelValues ...string) ScalableHistogramWithHeaders {
	next := make(MultiHistogramWithHeaders, len(h))
	for i := range h {
		next[i] = h[i].With(headers, labelValues...)
	}
	return next
}

// NewScalableHistogramWithNoopHeaders returns a multi-histogram, wrapping the passed ScalableHistogramWithHeaders.
func NewScalableHistogramWithNoopHeaders(h metrics.Histogram, unit time.Duration) (ScalableHistogramWithNoopHeaders, error) {
	sh, err := NewHistogramWithScale(h, unit)
	if err != nil {
		return ScalableHistogramWithNoopHeaders{}, err
	}
	return ScalableHistogramWithNoopHeaders{histogram: sh, unit: unit}, nil
}

// ScalableHistogramWithNoopHeaders collects multiple individual histograms and treats them as a unit.
type ScalableHistogramWithNoopHeaders struct {
	histogram ScalableHistogramWithHeaders
	unit      time.Duration
}

// ObserveFromStart implements ScalableHistogramWithHeaders.
func (h ScalableHistogramWithNoopHeaders) ObserveFromStart(start time.Time) {
	h.histogram.Observe(time.Since(start).Seconds())
}

// Observe implements ScalableHistogramWithHeaders.
func (h ScalableHistogramWithNoopHeaders) Observe(v float64) {
	h.histogram.Observe(v)
}

// With implements ScalableHistogramWithHeaders.
func (h ScalableHistogramWithNoopHeaders) With(headers http.Header, labelValues ...string) ScalableHistogramWithHeaders {
	return ScalableHistogramWithNoopHeaders{histogram: h.histogram.With(headers, labelValues...)}
}
