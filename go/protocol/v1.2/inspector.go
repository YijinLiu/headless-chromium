package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type InspectorEnableCB func(err error)

// Enables inspector domain notifications.
type InspectorEnableCommand struct {
	cb InspectorEnableCB
}

func NewInspectorEnableCommand(cb InspectorEnableCB) *InspectorEnableCommand {
	return &InspectorEnableCommand{
		cb: cb,
	}
}

func (cmd *InspectorEnableCommand) Name() string {
	return "Inspector.enable"
}

func (cmd *InspectorEnableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type InspectorDisableCB func(err error)

// Disables inspector domain notifications.
type InspectorDisableCommand struct {
	cb InspectorDisableCB
}

func NewInspectorDisableCommand(cb InspectorDisableCB) *InspectorDisableCommand {
	return &InspectorDisableCommand{
		cb: cb,
	}
}

func (cmd *InspectorDisableCommand) Name() string {
	return "Inspector.disable"
}

func (cmd *InspectorDisableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DetachedEvent struct {
	Reason string `json:"reason"` // The reason why connection has been terminated.
}

// Fired when remote debugging connection is about to be terminated. Contains detach reason.
type DetachedEventSink struct {
	events chan *DetachedEvent
}

func NewDetachedEventSink(bufSize int) *DetachedEventSink {
	return &DetachedEventSink{
		events: make(chan *DetachedEvent, bufSize),
	}
}

func (s *DetachedEventSink) Name() string {
	return "Inspector.detached"
}

func (s *DetachedEventSink) OnEvent(params []byte) {
	evt := &DetachedEvent{}
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

type TargetCrashedEvent struct {
}

// Fired when debugging target has crashed
type TargetCrashedEventSink struct {
	events chan *TargetCrashedEvent
}

func NewTargetCrashedEventSink(bufSize int) *TargetCrashedEventSink {
	return &TargetCrashedEventSink{
		events: make(chan *TargetCrashedEvent, bufSize),
	}
}

func (s *TargetCrashedEventSink) Name() string {
	return "Inspector.targetCrashed"
}

func (s *TargetCrashedEventSink) OnEvent(params []byte) {
	evt := &TargetCrashedEvent{}
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
