package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type BindParams struct {
	Port int `json:"port"` // Port number to bind.
}

// Request browser port binding.

type BindCommand struct {
	params *BindParams
	wg     sync.WaitGroup
	err    error
}

func NewBindCommand(params *BindParams) *BindCommand {
	return &BindCommand{
		params: params,
	}
}

func (cmd *BindCommand) Name() string {
	return "Tethering.bind"
}

func (cmd *BindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *BindCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Bind(params *BindParams, conn *hc.Conn) (err error) {
	cmd := NewBindCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type BindCB func(err error)

// Request browser port binding.

type AsyncBindCommand struct {
	params *BindParams
	cb     BindCB
}

func NewAsyncBindCommand(params *BindParams, cb BindCB) *AsyncBindCommand {
	return &AsyncBindCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncBindCommand) Name() string {
	return "Tethering.bind"
}

func (cmd *AsyncBindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *BindCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncBindCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type UnbindParams struct {
	Port int `json:"port"` // Port number to unbind.
}

// Request browser port unbinding.

type UnbindCommand struct {
	params *UnbindParams
	wg     sync.WaitGroup
	err    error
}

func NewUnbindCommand(params *UnbindParams) *UnbindCommand {
	return &UnbindCommand{
		params: params,
	}
}

func (cmd *UnbindCommand) Name() string {
	return "Tethering.unbind"
}

func (cmd *UnbindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnbindCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Unbind(params *UnbindParams, conn *hc.Conn) (err error) {
	cmd := NewUnbindCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type UnbindCB func(err error)

// Request browser port unbinding.

type AsyncUnbindCommand struct {
	params *UnbindParams
	cb     UnbindCB
}

func NewAsyncUnbindCommand(params *UnbindParams, cb UnbindCB) *AsyncUnbindCommand {
	return &AsyncUnbindCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncUnbindCommand) Name() string {
	return "Tethering.unbind"
}

func (cmd *AsyncUnbindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnbindCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncUnbindCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Informs that port was successfully bound and got a specified connection id.

type AcceptedEvent struct {
	Port         int    `json:"port"`         // Port number that was successfully bound.
	ConnectionId string `json:"connectionId"` // Connection id to be used.
}

func OnAccepted(conn *hc.Conn, cb func(evt *AcceptedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AcceptedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Tethering.accepted", sink)
}
