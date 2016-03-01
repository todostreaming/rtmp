package spec

import "io"

func ReadByte(r io.Reader) (b byte, err error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}

	buf := make([]byte, 1)
	if _, err = r.Read(buf); err != nil {
		return
	}

	b = buf[0]
	return
}

func ReadBytes(r io.Reader, n int) (buf []byte, err error) {
	buf = make([]byte, n)
	_, err = io.ReadFull(r, buf)

	return
}
