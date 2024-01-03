package conn_test

import (
	"testing"

	"github.com/todostreaming/rtmp/cmd/conn"
)

func TestRelevantCommandsImplementReceivable(t *testing.T) {
	for _, c := range []interface{}{
		new(conn.ConnectCommand),
		new(conn.CreateStreamCommand),
		new(conn.ReleaseCommand),
		new(conn.FCPublishCommand),
		new(conn.GetStreamLength),
	} {
		if _, receivable := c.(conn.Receivable); !receivable {
			t.Fatalf("type %T does not implement Receivable", c)
		}
	}
}
