package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type WorkerEnableCB func(err error)

type WorkerEnableCommand struct {
	cb WorkerEnableCB
}

func NewWorkerEnableCommand(cb WorkerEnableCB) *WorkerEnableCommand {
	return &WorkerEnableCommand{
		cb: cb,
	}
}

func (cmd *WorkerEnableCommand) Name() string {
	return "Worker.enable"
}

func (cmd *WorkerEnableCommand) Params() interface{} {
	return nil
}

func (cmd *WorkerEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type WorkerDisableCB func(err error)

type WorkerDisableCommand struct {
	cb WorkerDisableCB
}

func NewWorkerDisableCommand(cb WorkerDisableCB) *WorkerDisableCommand {
	return &WorkerDisableCommand{
		cb: cb,
	}
}

func (cmd *WorkerDisableCommand) Name() string {
	return "Worker.disable"
}

func (cmd *WorkerDisableCommand) Params() interface{} {
	return nil
}

func (cmd *WorkerDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SendMessageToWorkerParams struct {
	WorkerId string `json:"workerId"`
	Message  string `json:"message"`
}

type SendMessageToWorkerCB func(err error)

type SendMessageToWorkerCommand struct {
	params *SendMessageToWorkerParams
	cb     SendMessageToWorkerCB
}

func NewSendMessageToWorkerCommand(params *SendMessageToWorkerParams, cb SendMessageToWorkerCB) *SendMessageToWorkerCommand {
	return &SendMessageToWorkerCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SendMessageToWorkerCommand) Name() string {
	return "Worker.sendMessageToWorker"
}

func (cmd *SendMessageToWorkerCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SendMessageToWorkerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetWaitForDebuggerOnStartParams struct {
	Value bool `json:"value"`
}

type SetWaitForDebuggerOnStartCB func(err error)

type SetWaitForDebuggerOnStartCommand struct {
	params *SetWaitForDebuggerOnStartParams
	cb     SetWaitForDebuggerOnStartCB
}

func NewSetWaitForDebuggerOnStartCommand(params *SetWaitForDebuggerOnStartParams, cb SetWaitForDebuggerOnStartCB) *SetWaitForDebuggerOnStartCommand {
	return &SetWaitForDebuggerOnStartCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetWaitForDebuggerOnStartCommand) Name() string {
	return "Worker.setWaitForDebuggerOnStart"
}

func (cmd *SetWaitForDebuggerOnStartCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetWaitForDebuggerOnStartCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type WorkerCreatedEvent struct {
	WorkerId           string `json:"workerId"`
	Url                string `json:"url"`
	WaitingForDebugger bool   `json:"waitingForDebugger"`
}

type WorkerCreatedEventSink struct {
	events chan *WorkerCreatedEvent
}

func NewWorkerCreatedEventSink(bufSize int) *WorkerCreatedEventSink {
	return &WorkerCreatedEventSink{
		events: make(chan *WorkerCreatedEvent, bufSize),
	}
}

func (s *WorkerCreatedEventSink) Name() string {
	return "Worker.workerCreated"
}

func (s *WorkerCreatedEventSink) OnEvent(params []byte) {
	evt := &WorkerCreatedEvent{}
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

type WorkerTerminatedEvent struct {
	WorkerId string `json:"workerId"`
}

type WorkerTerminatedEventSink struct {
	events chan *WorkerTerminatedEvent
}

func NewWorkerTerminatedEventSink(bufSize int) *WorkerTerminatedEventSink {
	return &WorkerTerminatedEventSink{
		events: make(chan *WorkerTerminatedEvent, bufSize),
	}
}

func (s *WorkerTerminatedEventSink) Name() string {
	return "Worker.workerTerminated"
}

func (s *WorkerTerminatedEventSink) OnEvent(params []byte) {
	evt := &WorkerTerminatedEvent{}
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

type DispatchMessageFromWorkerEvent struct {
	WorkerId string `json:"workerId"`
	Message  string `json:"message"`
}

type DispatchMessageFromWorkerEventSink struct {
	events chan *DispatchMessageFromWorkerEvent
}

func NewDispatchMessageFromWorkerEventSink(bufSize int) *DispatchMessageFromWorkerEventSink {
	return &DispatchMessageFromWorkerEventSink{
		events: make(chan *DispatchMessageFromWorkerEvent, bufSize),
	}
}

func (s *DispatchMessageFromWorkerEventSink) Name() string {
	return "Worker.dispatchMessageFromWorker"
}

func (s *DispatchMessageFromWorkerEventSink) OnEvent(params []byte) {
	evt := &DispatchMessageFromWorkerEvent{}
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
