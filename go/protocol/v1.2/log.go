package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Log entry.
type LogEntry struct {
	Source           string            `json:"source"`           // Log entry source.
	Level            string            `json:"level"`            // Log entry severity.
	Text             string            `json:"text"`             // Logged text.
	Timestamp        *RuntimeTimestamp `json:"timestamp"`        // Timestamp when this entry was added.
	Url              string            `json:"url"`              // URL of the resource if known.
	LineNumber       int               `json:"lineNumber"`       // Line number in the resource.
	StackTrace       *StackTrace       `json:"stackTrace"`       // JavaScript stack trace.
	NetworkRequestId *RequestId        `json:"networkRequestId"` // Identifier of the network request associated with this entry.
	WorkerId         string            `json:"workerId"`         // Identifier of the worker associated with this entry.
}

type LogEnableCB func(err error)

// Enables log domain, sends the entries collected so far to the client by means of the entryAdded notification.
type LogEnableCommand struct {
	cb LogEnableCB
}

func NewLogEnableCommand(cb LogEnableCB) *LogEnableCommand {
	return &LogEnableCommand{
		cb: cb,
	}
}

func (cmd *LogEnableCommand) Name() string {
	return "Log.enable"
}

func (cmd *LogEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LogEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type LogDisableCB func(err error)

// Disables log domain, prevents further log entries from being reported to the client.
type LogDisableCommand struct {
	cb LogDisableCB
}

func NewLogDisableCommand(cb LogDisableCB) *LogDisableCommand {
	return &LogDisableCommand{
		cb: cb,
	}
}

func (cmd *LogDisableCommand) Name() string {
	return "Log.disable"
}

func (cmd *LogDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LogDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ClearCB func(err error)

// Clears the log.
type ClearCommand struct {
	cb ClearCB
}

func NewClearCommand(cb ClearCB) *ClearCommand {
	return &ClearCommand{
		cb: cb,
	}
}

func (cmd *ClearCommand) Name() string {
	return "Log.clear"
}

func (cmd *ClearCommand) Params() interface{} {
	return nil
}

func (cmd *ClearCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type EntryAddedEvent struct {
	Entry *LogEntry `json:"entry"` // The entry.
}

// Issued when new message was logged.
type EntryAddedEventSink struct {
	events chan *EntryAddedEvent
}

func NewEntryAddedEventSink(bufSize int) *EntryAddedEventSink {
	return &EntryAddedEventSink{
		events: make(chan *EntryAddedEvent, bufSize),
	}
}

func (s *EntryAddedEventSink) Name() string {
	return "Log.entryAdded"
}

func (s *EntryAddedEventSink) OnEvent(params []byte) {
	evt := &EntryAddedEvent{}
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
