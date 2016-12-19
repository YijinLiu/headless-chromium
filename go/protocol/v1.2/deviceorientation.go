package protocol

type DeviceOrientationSetDeviceOrientationOverrideParams struct {
	Alpha int `json:"alpha"` // Mock alpha
	Beta  int `json:"beta"`  // Mock beta
	Gamma int `json:"gamma"` // Mock gamma
}

type DeviceOrientationSetDeviceOrientationOverrideCB func(err error)

// Overrides the Device Orientation.
type DeviceOrientationSetDeviceOrientationOverrideCommand struct {
	params *DeviceOrientationSetDeviceOrientationOverrideParams
	cb     DeviceOrientationSetDeviceOrientationOverrideCB
}

func NewDeviceOrientationSetDeviceOrientationOverrideCommand(params *DeviceOrientationSetDeviceOrientationOverrideParams, cb DeviceOrientationSetDeviceOrientationOverrideCB) *DeviceOrientationSetDeviceOrientationOverrideCommand {
	return &DeviceOrientationSetDeviceOrientationOverrideCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.setDeviceOrientationOverride"
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Params() interface{} {
	return cmd.params
}

func (cmd *DeviceOrientationSetDeviceOrientationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DeviceOrientationClearDeviceOrientationOverrideCB func(err error)

// Clears the overridden Device Orientation.
type DeviceOrientationClearDeviceOrientationOverrideCommand struct {
	cb DeviceOrientationClearDeviceOrientationOverrideCB
}

func NewDeviceOrientationClearDeviceOrientationOverrideCommand(cb DeviceOrientationClearDeviceOrientationOverrideCB) *DeviceOrientationClearDeviceOrientationOverrideCommand {
	return &DeviceOrientationClearDeviceOrientationOverrideCommand{
		cb: cb,
	}
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Name() string {
	return "DeviceOrientation.clearDeviceOrientationOverride"
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Params() interface{} {
	return nil
}

func (cmd *DeviceOrientationClearDeviceOrientationOverrideCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
