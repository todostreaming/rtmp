package chunk

import "sync"

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
	// normalizer is the Normalizer type that produces complete chunks.
	normalizer Normalizer

	// smu guards streams
	smu sync.Mutex
	// streams maps chunk stream IDs (contained in the basic header of all
	// chunks) to their appropriate chunk Stream
	streams map[uint32]*Stream

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
func NewParser(reader Reader, normalizer Normalizer) *Parser {
	return &Parser{
		reader:     reader,
		normalizer: normalizer,
		streams:    make(map[uint32]*Stream),
		errs:       make(chan error),
		closer:     make(chan struct{}),
	}
}

// Stream returns the unique *chunk.Stream associated with a given chunk stream
// ID. Multiple calls to this method are guaranteed to return the same value. If
// a chunk stream does not exist for a given chunk stream ID, then one will be
// created (and subsequently returned). As a result, this method will never
// return nil.
func (p *Parser) Stream(id uint32) *Stream {
	p.smu.Lock()
	defer p.smu.Unlock()

	if _, ok := p.streams[id]; !ok {
		p.streams[id] = NewStream(id)
	}

	return p.streams[id]
}

// Errs returns a channel of errors which contains all reading errors
// encountered as a result of dealing with _any_ chunk stream.
func (p *Parser) Errs() <-chan error { return p.errs }

// Close halts the read/normalize process from all chunk streams and closes each
// "child" input channel of all `Stream`s.
func (p *Parser) Close() { p.closer <- struct{}{} }

// Recv is responsible for processing the chunks coming off of the underlying
// chunk.Reader. It first normalizes them and then places them onto the
// appropriate chunk stream, ensuring first that it exists. If an error is
// encountered, it is returned. If a close{} operation is sent, then the
// function will clean up after itself, and subsequently return.
//
// Recv runs within its own goroutine.
func (p *Parser) Recv() {
	go p.reader.Recv()

	for {
		select {
		case in := <-p.reader.Chunks():
			p.normalizer.Normalize(in)

			stream := p.Stream(in.StreamId())
			stream.in <- in
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
