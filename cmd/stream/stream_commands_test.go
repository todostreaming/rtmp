package stream_test

import (
	"testing"

	"github.com/todostreaming/rtmp/cmd/stream"
	"github.com/stretchr/testify/assert"
)

func TestStreamCommandsImplementCommand(t *testing.T) {
	for _, c := range []interface{}{
		new(stream.CommandPlay),
		new(stream.CommandPlay2),
		new(stream.CommandDeleteStream),
		new(stream.CommandReceiveAudio),
		new(stream.CommandReceiveVideo),
		new(stream.CommandPublish),
		new(stream.CommandSeek),
		new(stream.CommandPause),
	} {
		if cmd, ok := c.(stream.Command); ok {
			assert.True(t, cmd.IsCommand(),
				"expected type %T to truly implement command",
				cmd)
		} else {
			t.Errorf(
				"cmd/stream: type %T does not implement Command",
				c)
		}
	}
}
