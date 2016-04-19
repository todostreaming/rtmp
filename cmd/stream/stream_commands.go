package stream

import "github.com/WatchBeam/amf0"

type (
	CommandPlay struct {
		StreamName string
		Start      float64
		Duration   float64
		Reset      bool
	}

	CommandPlay2 struct {
		Parameters *amf0.Object
	}

	CommandDeleteStream struct {
		StreamId float64
	}

	CommandReceiveAudio struct {
		Successful bool
	}

	CommandReceiveVideo struct {
		Successful bool
	}

	CommandPublish struct {
		Name string
		Type string
	}

	CommandSeek struct {
		OffsetMillis float64
	}

	CommandPause struct {
		Paused       bool
		CutoffMillis float64
	}
)

func (_ *CommandPlay) IsCommand() bool         { return true }
func (_ *CommandPlay2) IsCommand() bool        { return true }
func (_ *CommandDeleteStream) IsCommand() bool { return true }
func (_ *CommandReceiveAudio) IsCommand() bool { return true }
func (_ *CommandReceiveVideo) IsCommand() bool { return true }
func (_ *CommandPublish) IsCommand() bool      { return true }
func (_ *CommandSeek) IsCommand() bool         { return true }
func (_ *CommandPause) IsCommand() bool        { return true }
