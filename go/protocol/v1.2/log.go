package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Log entry.
type LogEntry struct {
	Source           string            `json:"source"`                     // Log entry source.
	Level            string            `json:"level"`                      // Log entry severity.
	Text             string            `json:"text"`                       // Logged text.
	Timestamp        *RuntimeTimestamp `json:"timestamp"`                  // Timestamp when this entry was added.
	Url              string            `json:"url,omitempty"`              // URL of the resource if known.
	LineNumber       int               `json:"lineNumber,omitempty"`       // Line number in the resource.
	StackTrace       *StackTrace       `json:"stackTrace,omitempty"`       // JavaScript stack trace.
	NetworkRequestId *RequestId        `json:"networkRequestId,omitempty"` // Identifier of the network request associated with this entry.
	WorkerId         string            `json:"workerId,omitempty"`         // Identifier of the worker associated with this entry.
}

// Violation configuration setting.
type ViolationSetting struct {
	Name      string  `json:"name"`      // Violation type.
	Threshold float64 `json:"threshold"` // Time threshold to trigger upon.
}

// Enables log domain, sends the entries collected so far to the client by means of the entryAdded notification.

type LogEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewLogEnableCommand() *LogEnableCommand {
	return &LogEnableCommand{}
}

func (cmd *LogEnableCommand) Name() string {
	return "Log.enable"
}

func (cmd *LogEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LogEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func LogEnable(conn *hc.Conn) (err error) {
	cmd := NewLogEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type LogEnableCB func(err error)

// Enables log domain, sends the entries collected so far to the client by means of the entryAdded notification.

type AsyncLogEnableCommand struct {
	cb LogEnableCB
}

func NewAsyncLogEnableCommand(cb LogEnableCB) *AsyncLogEnableCommand {
	return &AsyncLogEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncLogEnableCommand) Name() string {
	return "Log.enable"
}

func (cmd *AsyncLogEnableCommand) Params() interface{} {
	return nil
}

func (cmd *LogEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncLogEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables log domain, prevents further log entries from being reported to the client.

type LogDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewLogDisableCommand() *LogDisableCommand {
	return &LogDisableCommand{}
}

func (cmd *LogDisableCommand) Name() string {
	return "Log.disable"
}

func (cmd *LogDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LogDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func LogDisable(conn *hc.Conn) (err error) {
	cmd := NewLogDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type LogDisableCB func(err error)

// Disables log domain, prevents further log entries from being reported to the client.

type AsyncLogDisableCommand struct {
	cb LogDisableCB
}

func NewAsyncLogDisableCommand(cb LogDisableCB) *AsyncLogDisableCommand {
	return &AsyncLogDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncLogDisableCommand) Name() string {
	return "Log.disable"
}

func (cmd *AsyncLogDisableCommand) Params() interface{} {
	return nil
}

func (cmd *LogDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncLogDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the log.

type ClearCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewClearCommand() *ClearCommand {
	return &ClearCommand{}
}

func (cmd *ClearCommand) Name() string {
	return "Log.clear"
}

func (cmd *ClearCommand) Params() interface{} {
	return nil
}

func (cmd *ClearCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func Clear(conn *hc.Conn) (err error) {
	cmd := NewClearCommand()
	cmd.Run(conn)
	return cmd.err
}

type ClearCB func(err error)

// Clears the log.

type AsyncClearCommand struct {
	cb ClearCB
}

func NewAsyncClearCommand(cb ClearCB) *AsyncClearCommand {
	return &AsyncClearCommand{
		cb: cb,
	}
}

func (cmd *AsyncClearCommand) Name() string {
	return "Log.clear"
}

func (cmd *AsyncClearCommand) Params() interface{} {
	return nil
}

func (cmd *ClearCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StartViolationsReportParams struct {
	Config []*ViolationSetting `json:"config"` // Configuration for violations.
}

// start violation reporting.

type StartViolationsReportCommand struct {
	params *StartViolationsReportParams
	wg     sync.WaitGroup
	err    error
}

func NewStartViolationsReportCommand(params *StartViolationsReportParams) *StartViolationsReportCommand {
	return &StartViolationsReportCommand{
		params: params,
	}
}

func (cmd *StartViolationsReportCommand) Name() string {
	return "Log.startViolationsReport"
}

func (cmd *StartViolationsReportCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartViolationsReportCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartViolationsReport(params *StartViolationsReportParams, conn *hc.Conn) (err error) {
	cmd := NewStartViolationsReportCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type StartViolationsReportCB func(err error)

// start violation reporting.

type AsyncStartViolationsReportCommand struct {
	params *StartViolationsReportParams
	cb     StartViolationsReportCB
}

func NewAsyncStartViolationsReportCommand(params *StartViolationsReportParams, cb StartViolationsReportCB) *AsyncStartViolationsReportCommand {
	return &AsyncStartViolationsReportCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncStartViolationsReportCommand) Name() string {
	return "Log.startViolationsReport"
}

func (cmd *AsyncStartViolationsReportCommand) Params() interface{} {
	return cmd.params
}

func (cmd *StartViolationsReportCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartViolationsReportCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Stop violation reporting.

type StopViolationsReportCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStopViolationsReportCommand() *StopViolationsReportCommand {
	return &StopViolationsReportCommand{}
}

func (cmd *StopViolationsReportCommand) Name() string {
	return "Log.stopViolationsReport"
}

func (cmd *StopViolationsReportCommand) Params() interface{} {
	return nil
}

func (cmd *StopViolationsReportCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopViolationsReport(conn *hc.Conn) (err error) {
	cmd := NewStopViolationsReportCommand()
	cmd.Run(conn)
	return cmd.err
}

type StopViolationsReportCB func(err error)

// Stop violation reporting.

type AsyncStopViolationsReportCommand struct {
	cb StopViolationsReportCB
}

func NewAsyncStopViolationsReportCommand(cb StopViolationsReportCB) *AsyncStopViolationsReportCommand {
	return &AsyncStopViolationsReportCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopViolationsReportCommand) Name() string {
	return "Log.stopViolationsReport"
}

func (cmd *AsyncStopViolationsReportCommand) Params() interface{} {
	return nil
}

func (cmd *StopViolationsReportCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopViolationsReportCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Issued when new message was logged.

type EntryAddedEvent struct {
	Entry *LogEntry `json:"entry"` // The entry.
}

func OnEntryAdded(conn *hc.Conn, cb func(evt *EntryAddedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &EntryAddedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Log.entryAdded", sink)
}
