package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Memory pressure level.
type PressureLevel string

const PressureLevelModerate PressureLevel = "moderate"
const PressureLevelCritical PressureLevel = "critical"

type GetDOMCountersResult struct {
	Documents        int `json:"documents"`
	Nodes            int `json:"nodes"`
	JsEventListeners int `json:"jsEventListeners"`
}

type GetDOMCountersCommand struct {
	result GetDOMCountersResult
	wg     sync.WaitGroup
	err    error
}

func NewGetDOMCountersCommand() *GetDOMCountersCommand {
	return &GetDOMCountersCommand{}
}

func (cmd *GetDOMCountersCommand) Name() string {
	return "Memory.getDOMCounters"
}

func (cmd *GetDOMCountersCommand) Params() interface{} {
	return nil
}

func (cmd *GetDOMCountersCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetDOMCounters(conn *hc.Conn) (result *GetDOMCountersResult, err error) {
	cmd := NewGetDOMCountersCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetDOMCountersCB func(result *GetDOMCountersResult, err error)

type AsyncGetDOMCountersCommand struct {
	cb GetDOMCountersCB
}

func NewAsyncGetDOMCountersCommand(cb GetDOMCountersCB) *AsyncGetDOMCountersCommand {
	return &AsyncGetDOMCountersCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetDOMCountersCommand) Name() string {
	return "Memory.getDOMCounters"
}

func (cmd *AsyncGetDOMCountersCommand) Params() interface{} {
	return nil
}

func (cmd *GetDOMCountersCommand) Result() *GetDOMCountersResult {
	return &cmd.result
}

func (cmd *GetDOMCountersCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetDOMCountersCommand) Done(data []byte, err error) {
	var result GetDOMCountersResult
	if err == nil {
		err = json.Unmarshal(data, &result)
	}
	if cmd.cb == nil {
		logging.Vlog(-1, err)
	} else if err != nil {
		cmd.cb(nil, err)
	} else {
		cmd.cb(&result, nil)
	}
}

type SetPressureNotificationsSuppressedParams struct {
	Suppressed bool `json:"suppressed"` // If true, memory pressure notifications will be suppressed.
}

// Enable/disable suppressing memory pressure notifications in all processes.

type SetPressureNotificationsSuppressedCommand struct {
	params *SetPressureNotificationsSuppressedParams
	wg     sync.WaitGroup
	err    error
}

func NewSetPressureNotificationsSuppressedCommand(params *SetPressureNotificationsSuppressedParams) *SetPressureNotificationsSuppressedCommand {
	return &SetPressureNotificationsSuppressedCommand{
		params: params,
	}
}

func (cmd *SetPressureNotificationsSuppressedCommand) Name() string {
	return "Memory.setPressureNotificationsSuppressed"
}

func (cmd *SetPressureNotificationsSuppressedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPressureNotificationsSuppressedCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetPressureNotificationsSuppressed(params *SetPressureNotificationsSuppressedParams, conn *hc.Conn) (err error) {
	cmd := NewSetPressureNotificationsSuppressedCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetPressureNotificationsSuppressedCB func(err error)

// Enable/disable suppressing memory pressure notifications in all processes.

type AsyncSetPressureNotificationsSuppressedCommand struct {
	params *SetPressureNotificationsSuppressedParams
	cb     SetPressureNotificationsSuppressedCB
}

func NewAsyncSetPressureNotificationsSuppressedCommand(params *SetPressureNotificationsSuppressedParams, cb SetPressureNotificationsSuppressedCB) *AsyncSetPressureNotificationsSuppressedCommand {
	return &AsyncSetPressureNotificationsSuppressedCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetPressureNotificationsSuppressedCommand) Name() string {
	return "Memory.setPressureNotificationsSuppressed"
}

func (cmd *AsyncSetPressureNotificationsSuppressedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPressureNotificationsSuppressedCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetPressureNotificationsSuppressedCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type SimulatePressureNotificationParams struct {
	Level PressureLevel `json:"level"` // Memory pressure level of the notification.
}

// Simulate a memory pressure notification in all processes.

type SimulatePressureNotificationCommand struct {
	params *SimulatePressureNotificationParams
	wg     sync.WaitGroup
	err    error
}

func NewSimulatePressureNotificationCommand(params *SimulatePressureNotificationParams) *SimulatePressureNotificationCommand {
	return &SimulatePressureNotificationCommand{
		params: params,
	}
}

func (cmd *SimulatePressureNotificationCommand) Name() string {
	return "Memory.simulatePressureNotification"
}

func (cmd *SimulatePressureNotificationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SimulatePressureNotificationCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SimulatePressureNotification(params *SimulatePressureNotificationParams, conn *hc.Conn) (err error) {
	cmd := NewSimulatePressureNotificationCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SimulatePressureNotificationCB func(err error)

// Simulate a memory pressure notification in all processes.

type AsyncSimulatePressureNotificationCommand struct {
	params *SimulatePressureNotificationParams
	cb     SimulatePressureNotificationCB
}

func NewAsyncSimulatePressureNotificationCommand(params *SimulatePressureNotificationParams, cb SimulatePressureNotificationCB) *AsyncSimulatePressureNotificationCommand {
	return &AsyncSimulatePressureNotificationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSimulatePressureNotificationCommand) Name() string {
	return "Memory.simulatePressureNotification"
}

func (cmd *AsyncSimulatePressureNotificationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SimulatePressureNotificationCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSimulatePressureNotificationCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
