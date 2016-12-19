package protocol

import (
	"encoding/json"
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

type GetDOMCountersCB func(result *GetDOMCountersResult, err error)

type GetDOMCountersCommand struct {
	cb GetDOMCountersCB
}

func NewGetDOMCountersCommand(cb GetDOMCountersCB) *GetDOMCountersCommand {
	return &GetDOMCountersCommand{
		cb: cb,
	}
}

func (cmd *GetDOMCountersCommand) Name() string {
	return "Memory.getDOMCounters"
}

func (cmd *GetDOMCountersCommand) Params() interface{} {
	return nil
}

func (cmd *GetDOMCountersCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetDOMCountersResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetPressureNotificationsSuppressedParams struct {
	Suppressed bool `json:"suppressed"` // If true, memory pressure notifications will be suppressed.
}

type SetPressureNotificationsSuppressedCB func(err error)

// Enable/disable suppressing memory pressure notifications in all processes.
type SetPressureNotificationsSuppressedCommand struct {
	params *SetPressureNotificationsSuppressedParams
	cb     SetPressureNotificationsSuppressedCB
}

func NewSetPressureNotificationsSuppressedCommand(params *SetPressureNotificationsSuppressedParams, cb SetPressureNotificationsSuppressedCB) *SetPressureNotificationsSuppressedCommand {
	return &SetPressureNotificationsSuppressedCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetPressureNotificationsSuppressedCommand) Name() string {
	return "Memory.setPressureNotificationsSuppressed"
}

func (cmd *SetPressureNotificationsSuppressedCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetPressureNotificationsSuppressedCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SimulatePressureNotificationParams struct {
	Level PressureLevel `json:"level"` // Memory pressure level of the notification.
}

type SimulatePressureNotificationCB func(err error)

// Simulate a memory pressure notification in all processes.
type SimulatePressureNotificationCommand struct {
	params *SimulatePressureNotificationParams
	cb     SimulatePressureNotificationCB
}

func NewSimulatePressureNotificationCommand(params *SimulatePressureNotificationParams, cb SimulatePressureNotificationCB) *SimulatePressureNotificationCommand {
	return &SimulatePressureNotificationCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SimulatePressureNotificationCommand) Name() string {
	return "Memory.simulatePressureNotification"
}

func (cmd *SimulatePressureNotificationCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SimulatePressureNotificationCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
