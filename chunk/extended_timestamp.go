package chunk

import (
	"io"

	"github.com/todostreaming/rtmp/spec"
)

// ExtendedTimestamp is an extra, optional, part of the standard RTMP chunk
// header. It is used to encode the complete 32-bit timestamp, or timestamp
// delta. It will only be present when the timestamp field of a type 0 header,
// or the timestamp delta field of a type 1 or 2 header is set to 0xffffff.
// Additionally, this field is present in type 3 chunks when the last type 0, 1
// or 2 chunk indicated presence of this field.
type ExtendedTimestamp struct {
	// Delta encodes the complete timestamp or extended timestamp field for
	// chunks matching the scenario as described above.
	Delta uint32
}

// Read slurps four bytes off of the given reader and parses the Delta field out
// as an unsigned, 32-bit integer.
func (t *ExtendedTimestamp) Read(r io.Reader) error {
	buf, err := spec.ReadBytes(r, 4)
	if err != nil {
		return err
	}

	t.Delta = spec.Uint32(buf)
	return nil
}

// Write serializes and writes the ExtendedTimestamp instance to the given
// io.Writer according to the RTMP specification.
func (t *ExtendedTimestamp) Write(w io.Writer) error {
	if _, err := spec.PutUint32(t.Delta, w); err != nil {
		return err
	}

	return nil
}
