package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Console message.
type ConsoleMessage struct {
	Source string `json:"source"`           // Message source.
	Level  string `json:"level"`            // Message severity.
	Text   string `json:"text"`             // Message text.
	Url    string `json:"url,omitempty"`    // URL of the message origin.
	Line   int    `json:"line,omitempty"`   // Line number in the resource that generated this message (1-based).
	Column int    `json:"column,omitempty"` // Column number in the resource that generated this message (1-based).
}

// Enables console domain, sends the messages collected so far to the client by means of the messageAdded notification.

type ConsoleEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewConsoleEnableCommand() *ConsoleEnableCommand {
	return &ConsoleEnableCommand{}
}

func (cmd *ConsoleEnableCommand) Name() string {
	return "Console.enable"
}

func (cmd *ConsoleEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ConsoleEnable(conn *hc.Conn) (err error) {
	cmd := NewConsoleEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ConsoleEnableCB func(err error)

// Enables console domain, sends the messages collected so far to the client by means of the messageAdded notification.

type AsyncConsoleEnableCommand struct {
	cb ConsoleEnableCB
}

func NewAsyncConsoleEnableCommand(cb ConsoleEnableCB) *AsyncConsoleEnableCommand {
	return &AsyncConsoleEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncConsoleEnableCommand) Name() string {
	return "Console.enable"
}

func (cmd *AsyncConsoleEnableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncConsoleEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables console domain, prevents further console messages from being reported to the client.

type ConsoleDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewConsoleDisableCommand() *ConsoleDisableCommand {
	return &ConsoleDisableCommand{}
}

func (cmd *ConsoleDisableCommand) Name() string {
	return "Console.disable"
}

func (cmd *ConsoleDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ConsoleDisable(conn *hc.Conn) (err error) {
	cmd := NewConsoleDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type ConsoleDisableCB func(err error)

// Disables console domain, prevents further console messages from being reported to the client.

type AsyncConsoleDisableCommand struct {
	cb ConsoleDisableCB
}

func NewAsyncConsoleDisableCommand(cb ConsoleDisableCB) *AsyncConsoleDisableCommand {
	return &AsyncConsoleDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncConsoleDisableCommand) Name() string {
	return "Console.disable"
}

func (cmd *AsyncConsoleDisableCommand) Params() interface{} {
	return nil
}

func (cmd *ConsoleDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncConsoleDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Does nothing.

type ClearMessagesCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewClearMessagesCommand() *ClearMessagesCommand {
	return &ClearMessagesCommand{}
}

func (cmd *ClearMessagesCommand) Name() string {
	return "Console.clearMessages"
}

func (cmd *ClearMessagesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearMessagesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ClearMessages(conn *hc.Conn) (err error) {
	cmd := NewClearMessagesCommand()
	cmd.Run(conn)
	return cmd.err
}

type ClearMessagesCB func(err error)

// Does nothing.

type AsyncClearMessagesCommand struct {
	cb ClearMessagesCB
}

func NewAsyncClearMessagesCommand(cb ClearMessagesCB) *AsyncClearMessagesCommand {
	return &AsyncClearMessagesCommand{
		cb: cb,
	}
}

func (cmd *AsyncClearMessagesCommand) Name() string {
	return "Console.clearMessages"
}

func (cmd *AsyncClearMessagesCommand) Params() interface{} {
	return nil
}

func (cmd *ClearMessagesCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearMessagesCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Issued when new console message is added.

type MessageAddedEvent struct {
	Message *ConsoleMessage `json:"message"` // Console message that has been added.
}

func OnMessageAdded(conn *hc.Conn, cb func(evt *MessageAddedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &MessageAddedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Console.messageAdded", sink)
}
