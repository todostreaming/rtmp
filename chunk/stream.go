package chunk

import (
	"encoding/binary"
	"io"
	"sync"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/WatchBeam/rtmp/spec"
)

type Stream struct {
	src io.Reader

	bmu      sync.Mutex
	builders map[uint32]*Builder

	rmu      sync.Mutex
	readSize int

	chunks chan *Chunk
	errs   chan error
}

func (s *Stream) Chunks() <-chan *Chunk { return s.chunks }
func (s *Stream) Errs() <-chan error    { return s.errs }

func (s *Stream) ReadSize() int {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	return s.readSize
}

func (s *Stream) SetReadSize(size int) {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	s.readSize = size
}

func (s *Stream) Recv() {
	for {
		header := new(Header)
		if err := header.Read(s.src); err != nil {
			s.errs <- err
			continue
		}

		builder := s.builder(header)
		n := spec.Min(builder.BytesLeft(), s.ReadSize())

		if err := builder.Read(s.src, n); err != nil {
			s.errs <- err
			continue
		}

		if builder.BytesLeft() == 0 {
			chunk := builder.Build()

			if !s.updateChunkSize(chunk) {
				s.chunks <- chunk
			}
			s.removeBuilder(chunk.Header.BasicHeader.StreamId)
		}
	}
}

func (s *Stream) updateChunkSize(chunk *Chunk) bool {
	if chunk.TypeId() != 0x01 {
		return false
	}

	s.SetReadSize(binary.BigEndian.Uint32(chunk.Data))
	return true
}

func (s *Stream) builder(header *Header) {
	s.bmu.Lock()
	defer s.bmu.Unlock()

	stream := header.BasicHeader.StreamId
	if s.builders[stream] == nil {
		s.builders[stream] = NewBuilder(header)
	}

	return s.builders[stream]
}

func (s *Stream) removeBuilder(stream uint32) {
	s.bmu.Lock()
	defer s.bmu.Unlock()

	delete(s.builders, chunk.Header.BasicHeader.StreamId)
}
