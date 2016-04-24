package chunk

import "sync"

// MultiStream represents a concatenation of multiple other chunk.Streams. It
// implements the Stream interface by providing an In() channel that is fed from
// all other `Append`-ed Streans by using the fan-in pattern.
type MultiStream struct {
	// wg waits on the number of active sources still open until that
	// number is zero.
	wg sync.WaitGroup
	// out is the sink of all active sources.
	out chan *Chunk
}

var _ Stream = new(MultiStream)

// NewMultiStream returns a new instance of the *MultiStream type with all
// internal channels initialized.
func NewMultiStream() *MultiStream {
	return &MultiStream{
		out: make(chan *Chunk),
	}
}

// In implements the chunk.Stream.In function. It provides an in-order read-only
// blocking channel which is the concatenation of all appended chunk streams.
func (m *MultiStream) In() <-chan *Chunk { return m.out }

// Append appends another chunk.Stream into this Multistream. It spawns another
// goroutine to `range` over the channel and feed it into the intenral out
// channel.
func (m *MultiStream) Append(others ...Stream) *MultiStream {
	m.wg.Add(len(others))

	for _, other := range others {
		go func(ch <-chan *Chunk) {
			defer m.wg.Done()

			for c := range ch {
				m.out <- c
			}
		}(other.In())
	}

	return m
}

// AwaitClose awaits until all feeding goroutines have closed themselves, and
// then closes our own internal out channel.
func (m *MultiStream) AwaitClose() {
	m.wg.Wait()

	close(m.out)
}
