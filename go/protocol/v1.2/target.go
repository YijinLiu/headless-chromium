package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type TargetTargetID string

type TargetType string

const TargetTypePage TargetType = "page"
const TargetTypeIframe TargetType = "iframe"
const TargetTypeWorker TargetType = "worker"
const TargetTypeService_worker TargetType = "service_worker"

type TargetTargetInfo struct {
	TargetId TargetTargetID `json:"targetId"`
	Type     TargetType     `json:"type"`
	Title    string         `json:"title"`
	Url      string         `json:"url"`
}

type SetDiscoverTargetsParams struct {
	Discover bool `json:"discover"` // Whether to discover available targets.
}

type SetDiscoverTargetsCB func(err error)

// Controls whether to discover available targets and notify via targetCreated/targetRemoved events.
type SetDiscoverTargetsCommand struct {
	params *SetDiscoverTargetsParams
	cb     SetDiscoverTargetsCB
}

func NewSetDiscoverTargetsCommand(params *SetDiscoverTargetsParams, cb SetDiscoverTargetsCB) *SetDiscoverTargetsCommand {
	return &SetDiscoverTargetsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetDiscoverTargetsCommand) Name() string {
	return "Target.setDiscoverTargets"
}

func (cmd *SetDiscoverTargetsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetDiscoverTargetsCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAutoAttachParams struct {
	AutoAttach             bool `json:"autoAttach"`             // Whether to auto-attach to related targets.
	WaitForDebuggerOnStart bool `json:"waitForDebuggerOnStart"` // Whether to pause new targets when attaching to them. Use Runtime.runIfWaitingForDebugger to run paused targets.
}

type SetAutoAttachCB func(err error)

// Controls whether to automatically attach to new targets which are considered to be related to this one. When turned on, attaches to all existing related targets as well. When turned off, automatically detaches from all currently attached targets.
type SetAutoAttachCommand struct {
	params *SetAutoAttachParams
	cb     SetAutoAttachCB
}

func NewSetAutoAttachCommand(params *SetAutoAttachParams, cb SetAutoAttachCB) *SetAutoAttachCommand {
	return &SetAutoAttachCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAutoAttachCommand) Name() string {
	return "Target.setAutoAttach"
}

func (cmd *SetAutoAttachCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAutoAttachCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SetAttachToFramesParams struct {
	Value bool `json:"value"` // Whether to attach to frames.
}

type SetAttachToFramesCB func(err error)

type SetAttachToFramesCommand struct {
	params *SetAttachToFramesParams
	cb     SetAttachToFramesCB
}

func NewSetAttachToFramesCommand(params *SetAttachToFramesParams, cb SetAttachToFramesCB) *SetAttachToFramesCommand {
	return &SetAttachToFramesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetAttachToFramesCommand) Name() string {
	return "Target.setAttachToFrames"
}

func (cmd *SetAttachToFramesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetAttachToFramesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SendMessageToTargetParams struct {
	TargetId string `json:"targetId"`
	Message  string `json:"message"`
}

type SendMessageToTargetCB func(err error)

type SendMessageToTargetCommand struct {
	params *SendMessageToTargetParams
	cb     SendMessageToTargetCB
}

func NewSendMessageToTargetCommand(params *SendMessageToTargetParams, cb SendMessageToTargetCB) *SendMessageToTargetCommand {
	return &SendMessageToTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SendMessageToTargetCommand) Name() string {
	return "Target.sendMessageToTarget"
}

func (cmd *SendMessageToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SendMessageToTargetCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetTargetInfoParams struct {
	TargetId TargetTargetID `json:"targetId"`
}

type GetTargetInfoResult struct {
	TargetInfo *TargetTargetInfo `json:"targetInfo"`
}

type GetTargetInfoCB func(result *GetTargetInfoResult, err error)

type GetTargetInfoCommand struct {
	params *GetTargetInfoParams
	cb     GetTargetInfoCB
}

func NewGetTargetInfoCommand(params *GetTargetInfoParams, cb GetTargetInfoCB) *GetTargetInfoCommand {
	return &GetTargetInfoCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetTargetInfoCommand) Name() string {
	return "Target.getTargetInfo"
}

func (cmd *GetTargetInfoCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetTargetInfoCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetTargetInfoResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ActivateTargetParams struct {
	TargetId TargetTargetID `json:"targetId"`
}

type ActivateTargetCB func(err error)

type ActivateTargetCommand struct {
	params *ActivateTargetParams
	cb     ActivateTargetCB
}

func NewActivateTargetCommand(params *ActivateTargetParams, cb ActivateTargetCB) *ActivateTargetCommand {
	return &ActivateTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ActivateTargetCommand) Name() string {
	return "Target.activateTarget"
}

func (cmd *ActivateTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ActivateTargetCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type AttachToTargetParams struct {
	TargetId TargetTargetID `json:"targetId"`
}

type AttachToTargetResult struct {
	Success bool `json:"success"` // Whether attach succeeded.
}

type AttachToTargetCB func(result *AttachToTargetResult, err error)

type AttachToTargetCommand struct {
	params *AttachToTargetParams
	cb     AttachToTargetCB
}

func NewAttachToTargetCommand(params *AttachToTargetParams, cb AttachToTargetCB) *AttachToTargetCommand {
	return &AttachToTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AttachToTargetCommand) Name() string {
	return "Target.attachToTarget"
}

func (cmd *AttachToTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AttachToTargetCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj AttachToTargetResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type DetachFromTargetParams struct {
	TargetId TargetTargetID `json:"targetId"`
}

type DetachFromTargetCB func(err error)

type DetachFromTargetCommand struct {
	params *DetachFromTargetParams
	cb     DetachFromTargetCB
}

func NewDetachFromTargetCommand(params *DetachFromTargetParams, cb DetachFromTargetCB) *DetachFromTargetCommand {
	return &DetachFromTargetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DetachFromTargetCommand) Name() string {
	return "Target.detachFromTarget"
}

func (cmd *DetachFromTargetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DetachFromTargetCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type TargetCreatedEvent struct {
	TargetInfo *TargetTargetInfo `json:"targetInfo"`
}

type TargetCreatedEventSink struct {
	events chan *TargetCreatedEvent
}

func NewTargetCreatedEventSink(bufSize int) *TargetCreatedEventSink {
	return &TargetCreatedEventSink{
		events: make(chan *TargetCreatedEvent, bufSize),
	}
}

func (s *TargetCreatedEventSink) Name() string {
	return "Target.targetCreated"
}

func (s *TargetCreatedEventSink) OnEvent(params []byte) {
	evt := &TargetCreatedEvent{}
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

type TargetRemovedEvent struct {
	TargetId TargetTargetID `json:"targetId"`
}

type TargetRemovedEventSink struct {
	events chan *TargetRemovedEvent
}

func NewTargetRemovedEventSink(bufSize int) *TargetRemovedEventSink {
	return &TargetRemovedEventSink{
		events: make(chan *TargetRemovedEvent, bufSize),
	}
}

func (s *TargetRemovedEventSink) Name() string {
	return "Target.targetRemoved"
}

func (s *TargetRemovedEventSink) OnEvent(params []byte) {
	evt := &TargetRemovedEvent{}
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

type AttachedToTargetEvent struct {
	TargetId           TargetTargetID `json:"targetId"`
	WaitingForDebugger bool           `json:"waitingForDebugger"`
}

type AttachedToTargetEventSink struct {
	events chan *AttachedToTargetEvent
}

func NewAttachedToTargetEventSink(bufSize int) *AttachedToTargetEventSink {
	return &AttachedToTargetEventSink{
		events: make(chan *AttachedToTargetEvent, bufSize),
	}
}

func (s *AttachedToTargetEventSink) Name() string {
	return "Target.attachedToTarget"
}

func (s *AttachedToTargetEventSink) OnEvent(params []byte) {
	evt := &AttachedToTargetEvent{}
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

type DetachedFromTargetEvent struct {
	TargetId TargetTargetID `json:"targetId"`
}

type DetachedFromTargetEventSink struct {
	events chan *DetachedFromTargetEvent
}

func NewDetachedFromTargetEventSink(bufSize int) *DetachedFromTargetEventSink {
	return &DetachedFromTargetEventSink{
		events: make(chan *DetachedFromTargetEvent, bufSize),
	}
}

func (s *DetachedFromTargetEventSink) Name() string {
	return "Target.detachedFromTarget"
}

func (s *DetachedFromTargetEventSink) OnEvent(params []byte) {
	evt := &DetachedFromTargetEvent{}
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

type ReceivedMessageFromTargetEvent struct {
	TargetId TargetTargetID `json:"targetId"`
	Message  string         `json:"message"`
}

type ReceivedMessageFromTargetEventSink struct {
	events chan *ReceivedMessageFromTargetEvent
}

func NewReceivedMessageFromTargetEventSink(bufSize int) *ReceivedMessageFromTargetEventSink {
	return &ReceivedMessageFromTargetEventSink{
		events: make(chan *ReceivedMessageFromTargetEvent, bufSize),
	}
}

func (s *ReceivedMessageFromTargetEventSink) Name() string {
	return "Target.receivedMessageFromTarget"
}

func (s *ReceivedMessageFromTargetEventSink) OnEvent(params []byte) {
	evt := &ReceivedMessageFromTargetEvent{}
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
