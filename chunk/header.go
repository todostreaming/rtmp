package chunk

type Header struct {
	BasicHeader       BasicHeader
	MessageHeader     MessageHeader
	ExtendedTimestamp ExtendedTimestamp
}
