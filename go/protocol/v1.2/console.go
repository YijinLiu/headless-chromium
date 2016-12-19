package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Console message.
type ConsoleMessage struct {
	Source string `json:"source"` // Message source.
	Level  string `json:"level"`  // Message severity.
	Text   string `json:"text"`   // Message text.
	Url    string `json:"url"`    // URL of the message origin.
	Line   int    `json:"line"`   // Line number in the resource that generated this message (1-based).
	Column int    `json:"column"` // Column number in the resource that generated this message (1-based).
}

type ConsoleEnableCB func(err error)

// Enables console domain, sends the messages collected so far to the client by means of the messageAdded notification.
type ConsoleEnableCommand struct {
	cb ConsoleEnableCB
}

func NewConsoleEnableCommand(cb ConsoleEnableCB) *ConsoleEnableCommand {
	return &ConsoleEnableCommand{
		cb: cb,
	}
}

func (cmd *ConsoleEnableCommand) Name() string {
	return "Console.enable"
}

func (cmd *ConsoleEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ConsoleDisableCB func(err error)

// Disables console domain, prevents further console messages from being reported to the client.
type ConsoleDisableCommand struct {
	cb ConsoleDisableCB
}

func NewConsoleDisableCommand(cb ConsoleDisableCB) *ConsoleDisableCommand {
	return &ConsoleDisableCommand{
		cb: cb,
	}
}

func (cmd *ConsoleDisableCommand) Name() string {
	return "Console.disable"
}

func (cmd *ConsoleDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ClearMessagesCB func(err error)

// Does nothing.
type ClearMessagesCommand struct {
	cb ClearMessagesCB
}

func NewClearMessagesCommand(cb ClearMessagesCB) *ClearMessagesCommand {
	return &ClearMessagesCommand{
		cb: cb,
	}
}

func (cmd *ClearMessagesCommand) Name() string {
	return "Console.clearMessages"
}

func (cmd *ClearMessagesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearMessagesCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type MessageAddedEvent struct {
	Message *ConsoleMessage `json:"message"` // Console message that has been added.
}

// Issued when new console message is added.
type MessageAddedEventSink struct {
	events chan *MessageAddedEvent
}

func NewMessageAddedEventSink(bufSize int) *MessageAddedEventSink {
	return &MessageAddedEventSink{
		events: make(chan *MessageAddedEvent, bufSize),
	}
}

func (s *MessageAddedEventSink) Name() string {
	return "Console.messageAdded"
}

func (s *MessageAddedEventSink) OnEvent(params []byte) {
	evt := &MessageAddedEvent{}
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
