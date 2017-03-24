package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
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
	Type   string  `json:"type"`             // Key type.
	Number float64 `json:"number,omitempty"` // Number value.
	String string  `json:"string,omitempty"` // String value.
	Date   float64 `json:"date,omitempty"`   // Date value.
	Array  []*Key  `json:"array,omitempty"`  // Array value.
}

// Key range.
type KeyRange struct {
	Lower     *Key `json:"lower,omitempty"` // Lower bound.
	Upper     *Key `json:"upper,omitempty"` // Upper bound.
	LowerOpen bool `json:"lowerOpen"`       // If true lower bound is open.
	UpperOpen bool `json:"upperOpen"`       // If true upper bound is open.
}

// Data entry.
type IndexedDBDataEntry struct {
	Key        *RemoteObject `json:"key"`        // Key object.
	PrimaryKey *RemoteObject `json:"primaryKey"` // Primary key object.
	Value      *RemoteObject `json:"value"`      // Value object.
}

// Key path.
type KeyPath struct {
	Type   string   `json:"type"`             // Key path type.
	String string   `json:"string,omitempty"` // String value.
	Array  []string `json:"array,omitempty"`  // Array value.
}

// Enables events from backend.

type IndexedDBEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewIndexedDBEnableCommand() *IndexedDBEnableCommand {
	return &IndexedDBEnableCommand{}
}

func (cmd *IndexedDBEnableCommand) Name() string {
	return "IndexedDB.enable"
}

func (cmd *IndexedDBEnableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func IndexedDBEnable(conn *hc.Conn) (err error) {
	cmd := NewIndexedDBEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type IndexedDBEnableCB func(err error)

// Enables events from backend.

type AsyncIndexedDBEnableCommand struct {
	cb IndexedDBEnableCB
}

func NewAsyncIndexedDBEnableCommand(cb IndexedDBEnableCB) *AsyncIndexedDBEnableCommand {
	return &AsyncIndexedDBEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncIndexedDBEnableCommand) Name() string {
	return "IndexedDB.enable"
}

func (cmd *AsyncIndexedDBEnableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncIndexedDBEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables events from backend.

type IndexedDBDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewIndexedDBDisableCommand() *IndexedDBDisableCommand {
	return &IndexedDBDisableCommand{}
}

func (cmd *IndexedDBDisableCommand) Name() string {
	return "IndexedDB.disable"
}

func (cmd *IndexedDBDisableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func IndexedDBDisable(conn *hc.Conn) (err error) {
	cmd := NewIndexedDBDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type IndexedDBDisableCB func(err error)

// Disables events from backend.

type AsyncIndexedDBDisableCommand struct {
	cb IndexedDBDisableCB
}

func NewAsyncIndexedDBDisableCommand(cb IndexedDBDisableCB) *AsyncIndexedDBDisableCommand {
	return &AsyncIndexedDBDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncIndexedDBDisableCommand) Name() string {
	return "IndexedDB.disable"
}

func (cmd *AsyncIndexedDBDisableCommand) Params() interface{} {
	return nil
}

func (cmd *IndexedDBDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncIndexedDBDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type RequestDatabaseNamesParams struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin.
}

type RequestDatabaseNamesResult struct {
	DatabaseNames []string `json:"databaseNames"` // Database names for origin.
}

// Requests database names for given security origin.

type RequestDatabaseNamesCommand struct {
	params *RequestDatabaseNamesParams
	result RequestDatabaseNamesResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestDatabaseNamesCommand(params *RequestDatabaseNamesParams) *RequestDatabaseNamesCommand {
	return &RequestDatabaseNamesCommand{
		params: params,
	}
}

func (cmd *RequestDatabaseNamesCommand) Name() string {
	return "IndexedDB.requestDatabaseNames"
}

func (cmd *RequestDatabaseNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseNamesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestDatabaseNames(params *RequestDatabaseNamesParams, conn *hc.Conn) (result *RequestDatabaseNamesResult, err error) {
	cmd := NewRequestDatabaseNamesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestDatabaseNamesCB func(result *RequestDatabaseNamesResult, err error)

// Requests database names for given security origin.

type AsyncRequestDatabaseNamesCommand struct {
	params *RequestDatabaseNamesParams
	cb     RequestDatabaseNamesCB
}

func NewAsyncRequestDatabaseNamesCommand(params *RequestDatabaseNamesParams, cb RequestDatabaseNamesCB) *AsyncRequestDatabaseNamesCommand {
	return &AsyncRequestDatabaseNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestDatabaseNamesCommand) Name() string {
	return "IndexedDB.requestDatabaseNames"
}

func (cmd *AsyncRequestDatabaseNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseNamesCommand) Result() *RequestDatabaseNamesResult {
	return &cmd.result
}

func (cmd *RequestDatabaseNamesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestDatabaseNamesCommand) Done(data []byte, err error) {
	var result RequestDatabaseNamesResult
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

type RequestDatabaseParams struct {
	SecurityOrigin string `json:"securityOrigin"` // Security origin.
	DatabaseName   string `json:"databaseName"`   // Database name.
}

type RequestDatabaseResult struct {
	DatabaseWithObjectStores *DatabaseWithObjectStores `json:"databaseWithObjectStores"` // Database with an array of object stores.
}

// Requests database with given name in given frame.

type RequestDatabaseCommand struct {
	params *RequestDatabaseParams
	result RequestDatabaseResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestDatabaseCommand(params *RequestDatabaseParams) *RequestDatabaseCommand {
	return &RequestDatabaseCommand{
		params: params,
	}
}

func (cmd *RequestDatabaseCommand) Name() string {
	return "IndexedDB.requestDatabase"
}

func (cmd *RequestDatabaseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestDatabase(params *RequestDatabaseParams, conn *hc.Conn) (result *RequestDatabaseResult, err error) {
	cmd := NewRequestDatabaseCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestDatabaseCB func(result *RequestDatabaseResult, err error)

// Requests database with given name in given frame.

type AsyncRequestDatabaseCommand struct {
	params *RequestDatabaseParams
	cb     RequestDatabaseCB
}

func NewAsyncRequestDatabaseCommand(params *RequestDatabaseParams, cb RequestDatabaseCB) *AsyncRequestDatabaseCommand {
	return &AsyncRequestDatabaseCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestDatabaseCommand) Name() string {
	return "IndexedDB.requestDatabase"
}

func (cmd *AsyncRequestDatabaseCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDatabaseCommand) Result() *RequestDatabaseResult {
	return &cmd.result
}

func (cmd *RequestDatabaseCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestDatabaseCommand) Done(data []byte, err error) {
	var result RequestDatabaseResult
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

type RequestDataParams struct {
	SecurityOrigin  string    `json:"securityOrigin"`     // Security origin.
	DatabaseName    string    `json:"databaseName"`       // Database name.
	ObjectStoreName string    `json:"objectStoreName"`    // Object store name.
	IndexName       string    `json:"indexName"`          // Index name, empty string for object store data requests.
	SkipCount       int       `json:"skipCount"`          // Number of records to skip.
	PageSize        int       `json:"pageSize"`           // Number of records to fetch.
	KeyRange        *KeyRange `json:"keyRange,omitempty"` // Key range.
}

type RequestDataResult struct {
	ObjectStoreDataEntries []*IndexedDBDataEntry `json:"objectStoreDataEntries"` // Array of object store data entries.
	HasMore                bool                  `json:"hasMore"`                // If true, there are more entries to fetch in the given range.
}

// Requests data from object store or index.

type RequestDataCommand struct {
	params *RequestDataParams
	result RequestDataResult
	wg     sync.WaitGroup
	err    error
}

func NewRequestDataCommand(params *RequestDataParams) *RequestDataCommand {
	return &RequestDataCommand{
		params: params,
	}
}

func (cmd *RequestDataCommand) Name() string {
	return "IndexedDB.requestData"
}

func (cmd *RequestDataCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDataCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func RequestData(params *RequestDataParams, conn *hc.Conn) (result *RequestDataResult, err error) {
	cmd := NewRequestDataCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type RequestDataCB func(result *RequestDataResult, err error)

// Requests data from object store or index.

type AsyncRequestDataCommand struct {
	params *RequestDataParams
	cb     RequestDataCB
}

func NewAsyncRequestDataCommand(params *RequestDataParams, cb RequestDataCB) *AsyncRequestDataCommand {
	return &AsyncRequestDataCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncRequestDataCommand) Name() string {
	return "IndexedDB.requestData"
}

func (cmd *AsyncRequestDataCommand) Params() interface{} {
	return cmd.params
}

func (cmd *RequestDataCommand) Result() *RequestDataResult {
	return &cmd.result
}

func (cmd *RequestDataCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncRequestDataCommand) Done(data []byte, err error) {
	var result RequestDataResult
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

type ClearObjectStoreParams struct {
	SecurityOrigin  string `json:"securityOrigin"`  // Security origin.
	DatabaseName    string `json:"databaseName"`    // Database name.
	ObjectStoreName string `json:"objectStoreName"` // Object store name.
}

// Clears all entries from an object store.

type ClearObjectStoreCommand struct {
	params *ClearObjectStoreParams
	wg     sync.WaitGroup
	err    error
}

func NewClearObjectStoreCommand(params *ClearObjectStoreParams) *ClearObjectStoreCommand {
	return &ClearObjectStoreCommand{
		params: params,
	}
}

func (cmd *ClearObjectStoreCommand) Name() string {
	return "IndexedDB.clearObjectStore"
}

func (cmd *ClearObjectStoreCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearObjectStoreCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ClearObjectStore(params *ClearObjectStoreParams, conn *hc.Conn) (err error) {
	cmd := NewClearObjectStoreCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ClearObjectStoreCB func(err error)

// Clears all entries from an object store.

type AsyncClearObjectStoreCommand struct {
	params *ClearObjectStoreParams
	cb     ClearObjectStoreCB
}

func NewAsyncClearObjectStoreCommand(params *ClearObjectStoreParams, cb ClearObjectStoreCB) *AsyncClearObjectStoreCommand {
	return &AsyncClearObjectStoreCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncClearObjectStoreCommand) Name() string {
	return "IndexedDB.clearObjectStore"
}

func (cmd *AsyncClearObjectStoreCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ClearObjectStoreCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncClearObjectStoreCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
