package protocol

import (
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type DeviceOrientationSetDeviceOrientationOverrideParams struct {
	Alpha float64 `json:"alpha"` // Mock alpha
	Beta  float64 `json:"beta"`  // Mock beta
	Gamma float64 `json:"gamma"` // Mock gamma
}

// Overrides the Device Orientation.

type DeviceOrientationSetDeviceOrientationOverrideCommand struct {
	params *DeviceOrientationSetDeviceOrientationOverrideParams
	wg     sync.WaitGroup
	err    error
}

func NewDeviceOrientationSetDeviceOrientationOverrideCommand(params *DeviceOrientationSetDeviceOrientationOverrideParams) *DeviceOrientationSetDeviceOrientationOverrideCommand {
	return &DeviceOrientationSetDeviceOrientationOverrideCommand{
		params: params,
	}
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.setDeviceOrientationOverride"
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DeviceOrientationSetDeviceOrientationOverride(params *DeviceOrientationSetDeviceOrientationOverrideParams, conn *hc.Conn) (err error) {
	cmd := NewDeviceOrientationSetDeviceOrientationOverrideCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type DeviceOrientationSetDeviceOrientationOverrideCB func(err error)

// Overrides the Device Orientation.

type AsyncDeviceOrientationSetDeviceOrientationOverrideCommand struct {
	params *DeviceOrientationSetDeviceOrientationOverrideParams
	cb     DeviceOrientationSetDeviceOrientationOverrideCB
}

func NewAsyncDeviceOrientationSetDeviceOrientationOverrideCommand(params *DeviceOrientationSetDeviceOrientationOverrideParams, cb DeviceOrientationSetDeviceOrientationOverrideCB) *AsyncDeviceOrientationSetDeviceOrientationOverrideCommand {
	return &AsyncDeviceOrientationSetDeviceOrientationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncDeviceOrientationSetDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.setDeviceOrientationOverride"
}

func (cmd *AsyncDeviceOrientationSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDeviceOrientationSetDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Clears the overridden Device Orientation.

type DeviceOrientationClearDeviceOrientationOverrideCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDeviceOrientationClearDeviceOrientationOverrideCommand() *DeviceOrientationClearDeviceOrientationOverrideCommand {
	return &DeviceOrientationClearDeviceOrientationOverrideCommand{}
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.clearDeviceOrientationOverride"
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DeviceOrientationClearDeviceOrientationOverride(conn *hc.Conn) (err error) {
	cmd := NewDeviceOrientationClearDeviceOrientationOverrideCommand()
	cmd.Run(conn)
	return cmd.err
}

type DeviceOrientationClearDeviceOrientationOverrideCB func(err error)

// Clears the overridden Device Orientation.

type AsyncDeviceOrientationClearDeviceOrientationOverrideCommand struct {
	cb DeviceOrientationClearDeviceOrientationOverrideCB
}

func NewAsyncDeviceOrientationClearDeviceOrientationOverrideCommand(cb DeviceOrientationClearDeviceOrientationOverrideCB) *AsyncDeviceOrientationClearDeviceOrientationOverrideCommand {
	return &AsyncDeviceOrientationClearDeviceOrientationOverrideCommand{
		cb: cb,
	}
}

func (cmd *AsyncDeviceOrientationClearDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.clearDeviceOrientationOverride"
}

func (cmd *AsyncDeviceOrientationClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDeviceOrientationClearDeviceOrientationOverrideCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
