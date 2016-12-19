package protocol

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

type ClearDataForOriginCB func(err error)

// Clears storage for origin.
type ClearDataForOriginCommand struct {
	params *ClearDataForOriginParams
	cb     ClearDataForOriginCB
}

func NewClearDataForOriginCommand(params *ClearDataForOriginParams, cb ClearDataForOriginCB) *ClearDataForOriginCommand {
	return &ClearDataForOriginCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ClearDataForOriginCommand) Name() string {
	return "Storage.clearDataForOrigin"
}

func (cmd *ClearDataForOriginCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearDataForOriginCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
