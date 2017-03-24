package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Enables inspector domain notifications.

type InspectorEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewInspectorEnableCommand() *InspectorEnableCommand {
	return &InspectorEnableCommand{}
}

func (cmd *InspectorEnableCommand) Name() string {
	return "Inspector.enable"
}

func (cmd *InspectorEnableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func InspectorEnable(conn *hc.Conn) (err error) {
	cmd := NewInspectorEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type InspectorEnableCB func(err error)

// Enables inspector domain notifications.

type AsyncInspectorEnableCommand struct {
	cb InspectorEnableCB
}

func NewAsyncInspectorEnableCommand(cb InspectorEnableCB) *AsyncInspectorEnableCommand {
	return &AsyncInspectorEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncInspectorEnableCommand) Name() string {
	return "Inspector.enable"
}

func (cmd *AsyncInspectorEnableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncInspectorEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables inspector domain notifications.

type InspectorDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewInspectorDisableCommand() *InspectorDisableCommand {
	return &InspectorDisableCommand{}
}

func (cmd *InspectorDisableCommand) Name() string {
	return "Inspector.disable"
}

func (cmd *InspectorDisableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func InspectorDisable(conn *hc.Conn) (err error) {
	cmd := NewInspectorDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type InspectorDisableCB func(err error)

// Disables inspector domain notifications.

type AsyncInspectorDisableCommand struct {
	cb InspectorDisableCB
}

func NewAsyncInspectorDisableCommand(cb InspectorDisableCB) *AsyncInspectorDisableCommand {
	return &AsyncInspectorDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncInspectorDisableCommand) Name() string {
	return "Inspector.disable"
}

func (cmd *AsyncInspectorDisableCommand) Params() interface{} {
	return nil
}

func (cmd *InspectorDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncInspectorDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Fired when remote debugging connection is about to be terminated. Contains detach reason.

type DetachedEvent struct {
	Reason string `json:"reason"` // The reason why connection has been terminated.
}

func OnDetached(conn *hc.Conn, cb func(evt *DetachedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &DetachedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Inspector.detached", sink)
}

// Fired when debugging target has crashed

type TargetCrashedEvent struct {
}

func OnTargetCrashed(conn *hc.Conn, cb func(evt *TargetCrashedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &TargetCrashedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Inspector.targetCrashed", sink)
}
