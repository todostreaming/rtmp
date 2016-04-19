package conn_test

import (
	"testing"

	"github.com/WatchBeam/rtmp/cmd/conn"
)

func TestRelevantResponsesImplementSendable(t *testing.T) {
	for _, c := range []interface{}{
		new(conn.CreateStreamResponse),
		new(conn.ConnectResponse),
	} {
		if _, sendable := c.(conn.Sendable); !sendable {
			t.Fatalf("type %T does not implement Sendable", c)
		}
	}
}
