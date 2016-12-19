package protocol

import (
	"encoding/json"
)

// Database with an array of object stores.
type DatabaseWithObjectStores struct {
	Name         string         `json:"name"`         // Database name.
	Version      int            `json:"version"`      // Database version.
	ObjectStores []*ObjectStore `json:"objectStores"` // Object stores in this database.
}

// Object store.
type ObjectStore struct {
	Name          string              `json:"name"`          // Object store name.
	KeyPath       *KeyPath            `json:"keyPath"`       // Object store key path.
	AutoIncrement bool                `json:"autoIncrement"` // If true, object store has auto increment flag set.
	Indexes       []*ObjectStoreIndex `json:"indexes"`       // Indexes in this object store.
}

// Object store index.
type ObjectStoreIndex struct {
	Name       string   `json:"name"`       // Index name.
	KeyPath    *KeyPath `json:"keyPath"`    // Index key path.
	Unique     bool     `json:"unique"`     // If true, index is unique.
	MultiEntry bool     `json:"multiEntry"` // If true, index allows multiple entries for a key.
}

// Key.
type Key struct {
	Type   string `json:"type"`   // Key type.
	Number int    `json:"number"` // Number value.
	String string `json:"string"` // String value.
	Date   int    `json:"date"`   // Date value.
	Array  []*Key `json:"array"`  // Array value.
}

// Key range.
type KeyRange struct {
	Lower     *Key `json:"lower"`     // Lower bound.
	Upper     *Key `json:"upper"`     // Upper bound.
	LowerOpen bool `json:"lowerOpen"` // If true lower bound is open.
	UpperOpen bool `json:"upperOpen"` // If true upper bound is open.
}

// Data entry.
type IndexedDBDataEntry struct {
	Key        *RemoteObject `json:"key"`        // Key object.
	PrimaryKey *RemoteObject `json:"primaryKey"` // Primary key object.
	Value      *RemoteObject `json:"value"`      // Value object.
}

// Key path.
type KeyPath struct {
	Type   string   `json:"type"`   // Key path type.
	String string   `json:"string"` // String value.
	Array  []string `json:"array"`  // Array value.
}

type IndexedDBEnableCB func(err error)

// Enables events from backend.
type IndexedDBEnableCommand struct {
	cb IndexedDBEnableCB
}

func NewIndexedDBEnableCommand(cb IndexedDBEnableCB) *IndexedDBEnableCommand {
	return &IndexedDBEnableCommand{
		cb: cb,
	}
}

func (cmd *IndexedDBEnableCommand) Name() string {
	return "IndexedDB.enable"
}

func (cmd *IndexedDBEnableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type IndexedDBDisableCB func(err error)

// Disables events from backend.
type IndexedDBDisableCommand struct {
	cb IndexedDBDisableCB
}

func NewIndexedDBDisableCommand(cb IndexedDBDisableCB) *IndexedDBDisableCommand {
	return &IndexedDBDisableCommand{
		cb: cb,
	}
}

func (cmd *IndexedDBDisableCommand) Name() string {
	return "IndexedDB.disable"
}

func (cmd *IndexedDBDisableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type RequestDatabaseNamesParams struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin.
}

type RequestDatabaseNamesResult struct {
	DatabaseNames []string `json:"databaseNames"` // Database names for origin.
}

type RequestDatabaseNamesCB func(result *RequestDatabaseNamesResult, err error)

// Requests database names for given security origin.
type RequestDatabaseNamesCommand struct {
	params *RequestDatabaseNamesParams
	cb     RequestDatabaseNamesCB
}

func NewRequestDatabaseNamesCommand(params *RequestDatabaseNamesParams, cb RequestDatabaseNamesCB) *RequestDatabaseNamesCommand {
	return &RequestDatabaseNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestDatabaseNamesCommand) Name() string {
	return "IndexedDB.requestDatabaseNames"
}

func (cmd *RequestDatabaseNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseNamesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestDatabaseNamesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RequestDatabaseParams struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin.
	DatabaseName   string `json:"databaseName"`   // Database name.
}

type RequestDatabaseResult struct {
	DatabaseWithObjectStores *DatabaseWithObjectStores `json:"databaseWithObjectStores"` // Database with an array of object stores.
}

type RequestDatabaseCB func(result *RequestDatabaseResult, err error)

// Requests database with given name in given frame.
type RequestDatabaseCommand struct {
	params *RequestDatabaseParams
	cb     RequestDatabaseCB
}

func NewRequestDatabaseCommand(params *RequestDatabaseParams, cb RequestDatabaseCB) *RequestDatabaseCommand {
	return &RequestDatabaseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestDatabaseCommand) Name() string {
	return "IndexedDB.requestDatabase"
}

func (cmd *RequestDatabaseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestDatabaseResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type RequestDataParams struct {
	SecurityOrigin  string    `json:"securityOrigin"`  // Security origin.
	DatabaseName    string    `json:"databaseName"`    // Database name.
	ObjectStoreName string    `json:"objectStoreName"` // Object store name.
	IndexName       string    `json:"indexName"`       // Index name, empty string for object store data requests.
	SkipCount       int       `json:"skipCount"`       // Number of records to skip.
	PageSize        int       `json:"pageSize"`        // Number of records to fetch.
	KeyRange        *KeyRange `json:"keyRange"`        // Key range.
}

type RequestDataResult struct {
	ObjectStoreDataEntries []*IndexedDBDataEntry `json:"objectStoreDataEntries"` // Array of object store data entries.
	HasMore                bool                  `json:"hasMore"`                // If true, there are more entries to fetch in the given range.
}

type RequestDataCB func(result *RequestDataResult, err error)

// Requests data from object store or index.
type RequestDataCommand struct {
	params *RequestDataParams
	cb     RequestDataCB
}

func NewRequestDataCommand(params *RequestDataParams, cb RequestDataCB) *RequestDataCommand {
	return &RequestDataCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *RequestDataCommand) Name() string {
	return "IndexedDB.requestData"
}

func (cmd *RequestDataCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDataCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj RequestDataResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ClearObjectStoreParams struct {
	SecurityOrigin  string `json:"securityOrigin"`  // Security origin.
	DatabaseName    string `json:"databaseName"`    // Database name.
	ObjectStoreName string `json:"objectStoreName"` // Object store name.
}

type ClearObjectStoreCB func(err error)

// Clears all entries from an object store.
type ClearObjectStoreCommand struct {
	params *ClearObjectStoreParams
	cb     ClearObjectStoreCB
}

func NewClearObjectStoreCommand(params *ClearObjectStoreParams, cb ClearObjectStoreCB) *ClearObjectStoreCommand {
	return &ClearObjectStoreCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ClearObjectStoreCommand) Name() string {
	return "IndexedDB.clearObjectStore"
}

func (cmd *ClearObjectStoreCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearObjectStoreCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}
