package chunk

import (
	"errors"
	"fmt"
	"sync"
)

// Parser is an intermediate chunk-parsing type that handles normalizing and
// splitting chunks received over the wire. Each chunk that is sent over _any_
// chunk stream is normalized, and then sent down the appropriate chunk stream.
//
// Chunks may be retreived in the following fashion:
//
// ```
// parser := chunk.NewParser()
// defer parser.Close()
//
//  #=> <chunk is sent over the network>
//  #=> <...>
//
// stream := parser.Stream(1)
// go func() {
//        for {
//                select {
//                case in := <-stream.In():
//                        /* ... */
//                case err := <-parser.Errs():
//                        /* ... */
//        }
// }()
// ```
type Parser struct {
	// reader is the Reader that chunks are read from.
	reader Reader

	// smu guards streams
	smu sync.Mutex
	// wg waits for the Recv loop to complete itself.
	wg sync.WaitGroup
	// streams maps chunk stream IDs (contained in the basic header of all
	// chunks) to their appropriate chunk Stream
	streams map[uint32]*stream

	// errs holds a channel of all errors encountered during the read/write
	// process.
	errs chan error
	// closer holds a channel that closes the Stream when anything is
	// written to it.
	closer chan struct{}
}

// NewParser allocates and returns a pointer to a new instance of the Parser
// type initialized with the given Reader and Normalizer.
//
// All internal channels, maps, and slices are instantiated at this time as
// well, and the receiving process is spawned in its own goroutine.
func NewParser(reader Reader) *Parser {
	return &Parser{
		reader:  reader,
		streams: make(map[uint32]*stream),
		errs:    make(chan error),
		closer:  make(chan struct{}),
	}
}

// Stream returns a chunk stream containing all of the IDs given as variadic
// arguments. This works in either one of two cases:
//
//  1) a single chunk stream (1 argument) is asked for, and either the one that
//  already exists, or a new instance of one is returned.
//
//  2) multiple chunk streams are asked for, and a MultiStream is returned
//  containing all of those chunk streams. If a single stream has already been
//  asked for in the set of streams to concatenate, an error is returned, and no
//  new chunk streams are created.
func (p *Parser) Stream(ids ...uint32) (Stream, error) {
	if len(ids) == 0 {
		return nil, errors.New(
			"rtmp/chunk: cannot return empty chunk stream")
	}

	p.smu.Lock()
	defer p.smu.Unlock()

	if len(ids) == 1 {
		id := ids[0]

		if _, ok := p.streams[id]; !ok {
			p.streams[id] = NewStream(id)
		}

		return p.streams[id], nil
	}

	for _, id := range ids {
		if _, exists := p.streams[id]; exists {
			return nil, fmt.Errorf(
				"rtmp/chunk: stream %v already exists", id)
		}
	}

	multi := NewMultiStream()

	for _, id := range ids {
		stream := NewStream(id)
		p.streams[id] = stream

		multi.Append(stream)
	}

	return multi, nil
}

// Errs returns a channel of errors which contains all reading errors
// encountered as a result of dealing with _any_ chunk stream.
func (p *Parser) Errs() <-chan error { return p.errs }

// Close halts the read/normalize process from all chunk streams and closes each
// "child" input channel of all `Stream`s.
func (p *Parser) Close() { p.closer <- struct{}{}; p.wg.Wait() }

// Recv is responsible for processing the chunks coming off of the underlying
// chunk.Reader. It first normalizes them and then places them onto the
// appropriate chunk stream, ensuring first that it exists. If an error is
// encountered, it is returned. If a close{} operation is sent, then the
// function will clean up after itself, and subsequently return.
//
// Recv runs within its own goroutine.
func (p *Parser) Recv() {
	p.wg.Add(1)
	defer p.wg.Done()

	go p.reader.Recv()

	for {
		select {
		case in := <-p.reader.Chunks():
			s, err := p.Stream(in.StreamId())
			if err != nil {
				p.errs <- err
				continue
			}

			s.(*stream).in <- in
		case err := <-p.reader.Errs():
			p.errs <- err
		case <-p.closer:
			p.reader.Close()

			close(p.errs)
			close(p.closer)

			p.smu.Lock()
			for _, stream := range p.streams {
				close(stream.in)
			}
			p.smu.Unlock()

			return
		}
	}
}
