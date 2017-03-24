package protocol

import (
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Enum of possible storage types.
type StorageType string

const StorageTypeAppcache StorageType = "appcache"
const StorageTypeCookies StorageType = "cookies"
const StorageTypeFile_systems StorageType = "file_systems"
const StorageTypeIndexeddb StorageType = "indexeddb"
const StorageTypeLocal_storage StorageType = "local_storage"
const StorageTypeShader_cache StorageType = "shader_cache"
const StorageTypeWebsql StorageType = "websql"
const StorageTypeService_workers StorageType = "service_workers"
const StorageTypeCache_storage StorageType = "cache_storage"
const StorageTypeAll StorageType = "all"

type ClearDataForOriginParams struct {
	Origin       string `json:"origin"`       // Security origin.
	StorageTypes string `json:"storageTypes"` // Comma separated origin names.
}

// Clears storage for origin.

type ClearDataForOriginCommand struct {
	params *ClearDataForOriginParams
	wg     sync.WaitGroup
	err    error
}

func NewClearDataForOriginCommand(params *ClearDataForOriginParams) *ClearDataForOriginCommand {
	return &ClearDataForOriginCommand{
		params: params,
	}
}

func (cmd *ClearDataForOriginCommand) Name() string {
	return "Storage.clearDataForOrigin"
}

func (cmd *ClearDataForOriginCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearDataForOriginCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ClearDataForOrigin(params *ClearDataForOriginParams, conn *hc.Conn) (err error) {
	cmd := NewClearDataForOriginCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ClearDataForOriginCB func(err error)

// Clears storage for origin.

type AsyncClearDataForOriginCommand struct {
	params *ClearDataForOriginParams
	cb     ClearDataForOriginCB
}

func NewAsyncClearDataForOriginCommand(params *ClearDataForOriginParams, cb ClearDataForOriginCB) *AsyncClearDataForOriginCommand {
	return &AsyncClearDataForOriginCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncClearDataForOriginCommand) Name() string {
	return "Storage.clearDataForOrigin"
}

func (cmd *AsyncClearDataForOriginCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearDataForOriginCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearDataForOriginCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
