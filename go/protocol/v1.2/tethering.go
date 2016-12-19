package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type BindParams struct {
	Port int `json:"port"` // Port number to bind.
}

type BindCB func(err error)

// Request browser port binding.
type BindCommand struct {
	params *BindParams
	cb     BindCB
}

func NewBindCommand(params *BindParams, cb BindCB) *BindCommand {
	return &BindCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *BindCommand) Name() string {
	return "Tethering.bind"
}

func (cmd *BindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *BindCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type UnbindParams struct {
	Port int `json:"port"` // Port number to unbind.
}

type UnbindCB func(err error)

// Request browser port unbinding.
type UnbindCommand struct {
	params *UnbindParams
	cb     UnbindCB
}

func NewUnbindCommand(params *UnbindParams, cb UnbindCB) *UnbindCommand {
	return &UnbindCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *UnbindCommand) Name() string {
	return "Tethering.unbind"
}

func (cmd *UnbindCommand) Params() interface{} {
	return cmd.params
}

func (cmd *UnbindCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type AcceptedEvent struct {
	Port         int    `json:"port"`         // Port number that was successfully bound.
	ConnectionId string `json:"connectionId"` // Connection id to be used.
}

// Informs that port was successfully bound and got a specified connection id.
type AcceptedEventSink struct {
	events chan *AcceptedEvent
}

func NewAcceptedEventSink(bufSize int) *AcceptedEventSink {
	return &AcceptedEventSink{
		events: make(chan *AcceptedEvent, bufSize),
	}
}

func (s *AcceptedEventSink) Name() string {
	return "Tethering.accepted"
}

func (s *AcceptedEventSink) OnEvent(params []byte) {
	evt := &AcceptedEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}
