package protocol

import (
	"encoding/json"
)

type StreamHandle string

type ReadParams struct {
	Handle StreamHandle `json:"handle"` // Handle of the stream to read.
	Offset int          `json:"offset"` // Seek to the specified offset before reading (if not specificed, proceed with offset following the last read).
	Size   int          `json:"size"`   // Maximum number of bytes to read (left upon the agent discretion if not specified).
}

type ReadResult struct {
	Data string `json:"data"` // Data that were read.
	Eof  bool   `json:"eof"`  // Set if the end-of-file condition occured while reading.
}

type ReadCB func(result *ReadResult, err error)

// Read a chunk of the stream
type ReadCommand struct {
	params *ReadParams
	cb     ReadCB
}

func NewReadCommand(params *ReadParams, cb ReadCB) *ReadCommand {
	return &ReadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ReadCommand) Name() string {
	return "IO.read"
}

func (cmd *ReadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReadCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ReadResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CloseParams struct {
	Handle StreamHandle `json:"handle"` // Handle of the stream to close.
}

type CloseCB func(err error)

// Close the stream, discard any temporary backing storage.
type CloseCommand struct {
	params *CloseParams
	cb     CloseCB
}

func NewCloseCommand(params *CloseParams, cb CloseCB) *CloseCommand {
	return &CloseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CloseCommand) Name() string {
	return "IO.close"
}

func (cmd *CloseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
