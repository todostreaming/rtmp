package chunk_test

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/stretchr/testify/assert"
)

func TestExtendedTimestampReadsCorrectly(t *testing.T) {
	time := rand.Uint32()
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, time)

	e := &chunk.ExtendedTimestamp{}
	err := e.Read(bytes.NewBuffer(buf))

	assert.Nil(t, err)
	assert.Equal(t, time, e.Delta)
}
