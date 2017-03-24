package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Describes a single graphics processor (GPU).
type GPUDevice struct {
	VendorId     float64 `json:"vendorId"`     // PCI ID of the GPU vendor, if available; 0 otherwise.
	DeviceId     float64 `json:"deviceId"`     // PCI ID of the GPU device, if available; 0 otherwise.
	VendorString string  `json:"vendorString"` // String description of the GPU vendor, if the PCI ID is not available.
	DeviceString string  `json:"deviceString"` // String description of the GPU device, if the PCI ID is not available.
}

// Provides information about the GPU(s) on the system.
type GPUInfo struct {
	Devices              []*GPUDevice      `json:"devices"`                 // The graphics devices on the system. Element 0 is the primary GPU.
	AuxAttributes        map[string]string `json:"auxAttributes,omitempty"` // An optional dictionary of additional GPU related attributes.
	FeatureStatus        map[string]string `json:"featureStatus,omitempty"` // An optional dictionary of graphics features and their status.
	DriverBugWorkarounds []string          `json:"driverBugWorkarounds"`    // An optional array of GPU driver bug workarounds.
}

type GetInfoResult struct {
	Gpu          *GPUInfo `json:"gpu"`          // Information about the GPUs on the system.
	ModelName    string   `json:"modelName"`    // A platform-dependent description of the model of the machine. On Mac OS, this is, for example, 'MacBookPro'. Will be the empty string if not supported.
	ModelVersion string   `json:"modelVersion"` // A platform-dependent description of the version of the machine. On Mac OS, this is, for example, '10.1'. Will be the empty string if not supported.
}

// Returns information about the system.

type GetInfoCommand struct {
	result GetInfoResult
	wg     sync.WaitGroup
	err    error
}

func NewGetInfoCommand() *GetInfoCommand {
	return &GetInfoCommand{}
}

func (cmd *GetInfoCommand) Name() string {
	return "SystemInfo.getInfo"
}

func (cmd *GetInfoCommand) Params() interface{} {
	return nil
}

func (cmd *GetInfoCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetInfo(conn *hc.Conn) (result *GetInfoResult, err error) {
	cmd := NewGetInfoCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetInfoCB func(result *GetInfoResult, err error)

// Returns information about the system.

type AsyncGetInfoCommand struct {
	cb GetInfoCB
}

func NewAsyncGetInfoCommand(cb GetInfoCB) *AsyncGetInfoCommand {
	return &AsyncGetInfoCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetInfoCommand) Name() string {
	return "SystemInfo.getInfo"
}

func (cmd *AsyncGetInfoCommand) Params() interface{} {
	return nil
}

func (cmd *GetInfoCommand) Result() *GetInfoResult {
	return &cmd.result
}

func (cmd *GetInfoCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetInfoCommand) Done(data []byte, err error) {
	var result GetInfoResult
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
