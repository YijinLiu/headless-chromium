package protocol

import (
	"encoding/json"
)

// Describes a single graphics processor (GPU).
type GPUDevice struct {
	VendorId     int    `json:"vendorId"`     // PCI ID of the GPU vendor, if available; 0 otherwise.
	DeviceId     int    `json:"deviceId"`     // PCI ID of the GPU device, if available; 0 otherwise.
	VendorString string `json:"vendorString"` // String description of the GPU vendor, if the PCI ID is not available.
	DeviceString string `json:"deviceString"` // String description of the GPU device, if the PCI ID is not available.
}

// Provides information about the GPU(s) on the system.
type GPUInfo struct {
	Devices              []*GPUDevice      `json:"devices"`              // The graphics devices on the system. Element 0 is the primary GPU.
	AuxAttributes        map[string]string `json:"auxAttributes"`        // An optional dictionary of additional GPU related attributes.
	FeatureStatus        map[string]string `json:"featureStatus"`        // An optional dictionary of graphics features and their status.
	DriverBugWorkarounds []string          `json:"driverBugWorkarounds"` // An optional array of GPU driver bug workarounds.
}

type GetInfoResult struct {
	Gpu          *GPUInfo `json:"gpu"`          // Information about the GPUs on the system.
	ModelName    string   `json:"modelName"`    // A platform-dependent description of the model of the machine. On Mac OS, this is, for example, 'MacBookPro'. Will be the empty string if not supported.
	ModelVersion string   `json:"modelVersion"` // A platform-dependent description of the version of the machine. On Mac OS, this is, for example, '10.1'. Will be the empty string if not supported.
}

type GetInfoCB func(result *GetInfoResult, err error)

// Returns information about the system.
type GetInfoCommand struct {
	cb GetInfoCB
}

func NewGetInfoCommand(cb GetInfoCB) *GetInfoCommand {
	return &GetInfoCommand{
		cb: cb,
	}
}

func (cmd *GetInfoCommand) Name() string {
	return "SystemInfo.getInfo"
}

func (cmd *GetInfoCommand) Params() interface{} {
	return nil
}

func (cmd *GetInfoCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetInfoResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}
