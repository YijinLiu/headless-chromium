package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type StreamHandle string

type ReadParams struct {
	Handle StreamHandle `json:"handle"`           // Handle of the stream to read.
	Offset int          `json:"offset,omitempty"` // Seek to the specified offset before reading (if not specificed, proceed with offset following the last read).
	Size   int          `json:"size,omitempty"`   // Maximum number of bytes to read (left upon the agent discretion if not specified).
}

type ReadResult struct {
	Data string `json:"data"` // Data that were read.
	Eof  bool   `json:"eof"`  // Set if the end-of-file condition occured while reading.
}

// Read a chunk of the stream

type ReadCommand struct {
	params *ReadParams
	result ReadResult
	wg     sync.WaitGroup
	err    error
}

func NewReadCommand(params *ReadParams) *ReadCommand {
	return &ReadCommand{
		params: params,
	}
}

func (cmd *ReadCommand) Name() string {
	return "IO.read"
}

func (cmd *ReadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReadCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Read(params *ReadParams, conn *hc.Conn) (result *ReadResult, err error) {
	cmd := NewReadCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ReadCB func(result *ReadResult, err error)

// Read a chunk of the stream

type AsyncReadCommand struct {
	params *ReadParams
	cb     ReadCB
}

func NewAsyncReadCommand(params *ReadParams, cb ReadCB) *AsyncReadCommand {
	return &AsyncReadCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncReadCommand) Name() string {
	return "IO.read"
}

func (cmd *AsyncReadCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ReadCommand) Result() *ReadResult {
	return &cmd.result
}

func (cmd *ReadCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncReadCommand) Done(data []byte, err error) {
	var result ReadResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type CloseParams struct {
	Handle StreamHandle `json:"handle"` // Handle of the stream to close.
}

// Close the stream, discard any temporary backing storage.

type CloseCommand struct {
	params *CloseParams
	wg     sync.WaitGroup
	err    error
}

func NewCloseCommand(params *CloseParams) *CloseCommand {
	return &CloseCommand{
		params: params,
	}
}

func (cmd *CloseCommand) Name() string {
	return "IO.close"
}

func (cmd *CloseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Close(params *CloseParams, conn *hc.Conn) (err error) {
	cmd := NewCloseCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type CloseCB func(err error)

// Close the stream, discard any temporary backing storage.

type AsyncCloseCommand struct {
	params *CloseParams
	cb     CloseCB
}

func NewAsyncCloseCommand(params *CloseParams, cb CloseCB) *AsyncCloseCommand {
	return &AsyncCloseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCloseCommand) Name() string {
	return "IO.close"
}

func (cmd *AsyncCloseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CloseCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCloseCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
