package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Unique identifier of Database object.
// @experimental
type DatabaseId string

// Database object.
// @experimental
type Database struct {
	Id      DatabaseId `json:"id"`      // Database ID.
	Domain  string     `json:"domain"`  // Database domain.
	Name    string     `json:"name"`    // Database name.
	Version string     `json:"version"` // Database version.
}

// Database error.
type Error struct {
	Message string `json:"message"` // Error message.
	Code    int    `json:"code"`    // Error code.
}

// Enables database tracking, database events will now be delivered to the client.

type DatabaseEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDatabaseEnableCommand() *DatabaseEnableCommand {
	return &DatabaseEnableCommand{}
}

func (cmd *DatabaseEnableCommand) Name() string {
	return "Database.enable"
}

func (cmd *DatabaseEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DatabaseEnable(conn *hc.Conn) (err error) {
	cmd := NewDatabaseEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DatabaseEnableCB func(err error)

// Enables database tracking, database events will now be delivered to the client.

type AsyncDatabaseEnableCommand struct {
	cb DatabaseEnableCB
}

func NewAsyncDatabaseEnableCommand(cb DatabaseEnableCB) *AsyncDatabaseEnableCommand {
	return &AsyncDatabaseEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDatabaseEnableCommand) Name() string {
	return "Database.enable"
}

func (cmd *AsyncDatabaseEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDatabaseEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables database tracking, prevents database events from being sent to the client.

type DatabaseDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewDatabaseDisableCommand() *DatabaseDisableCommand {
	return &DatabaseDisableCommand{}
}

func (cmd *DatabaseDisableCommand) Name() string {
	return "Database.disable"
}

func (cmd *DatabaseDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func DatabaseDisable(conn *hc.Conn) (err error) {
	cmd := NewDatabaseDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type DatabaseDisableCB func(err error)

// Disables database tracking, prevents database events from being sent to the client.

type AsyncDatabaseDisableCommand struct {
	cb DatabaseDisableCB
}

func NewAsyncDatabaseDisableCommand(cb DatabaseDisableCB) *AsyncDatabaseDisableCommand {
	return &AsyncDatabaseDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncDatabaseDisableCommand) Name() string {
	return "Database.disable"
}

func (cmd *AsyncDatabaseDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncDatabaseDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetDatabaseTableNamesParams struct {
	DatabaseId DatabaseId `json:"databaseId"`
}

type GetDatabaseTableNamesResult struct {
	TableNames []string `json:"tableNames"`
}

type GetDatabaseTableNamesCommand struct {
	params *GetDatabaseTableNamesParams
	result GetDatabaseTableNamesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetDatabaseTableNamesCommand(params *GetDatabaseTableNamesParams) *GetDatabaseTableNamesCommand {
	return &GetDatabaseTableNamesCommand{
		params: params,
	}
}

func (cmd *GetDatabaseTableNamesCommand) Name() string {
	return "Database.getDatabaseTableNames"
}

func (cmd *GetDatabaseTableNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDatabaseTableNamesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetDatabaseTableNames(params *GetDatabaseTableNamesParams, conn *hc.Conn) (result *GetDatabaseTableNamesResult, err error) {
	cmd := NewGetDatabaseTableNamesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetDatabaseTableNamesCB func(result *GetDatabaseTableNamesResult, err error)

type AsyncGetDatabaseTableNamesCommand struct {
	params *GetDatabaseTableNamesParams
	cb     GetDatabaseTableNamesCB
}

func NewAsyncGetDatabaseTableNamesCommand(params *GetDatabaseTableNamesParams, cb GetDatabaseTableNamesCB) *AsyncGetDatabaseTableNamesCommand {
	return &AsyncGetDatabaseTableNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetDatabaseTableNamesCommand) Name() string {
	return "Database.getDatabaseTableNames"
}

func (cmd *AsyncGetDatabaseTableNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDatabaseTableNamesCommand) Result() *GetDatabaseTableNamesResult {
	return &cmd.result
}

func (cmd *GetDatabaseTableNamesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetDatabaseTableNamesCommand) Done(data []byte, err error) {
	var result GetDatabaseTableNamesResult
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

type ExecuteSQLParams struct {
	DatabaseId DatabaseId `json:"databaseId"`
	Query      string     `json:"query"`
}

type ExecuteSQLResult struct {
	ColumnNames []string          `json:"columnNames"`
	Values      []json.RawMessage `json:"values"`
	SqlError    *Error            `json:"sqlError"`
}

type ExecuteSQLCommand struct {
	params *ExecuteSQLParams
	result ExecuteSQLResult
	wg     sync.WaitGroup
	err    error
}

func NewExecuteSQLCommand(params *ExecuteSQLParams) *ExecuteSQLCommand {
	return &ExecuteSQLCommand{
		params: params,
	}
}

func (cmd *ExecuteSQLCommand) Name() string {
	return "Database.executeSQL"
}

func (cmd *ExecuteSQLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ExecuteSQLCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ExecuteSQL(params *ExecuteSQLParams, conn *hc.Conn) (result *ExecuteSQLResult, err error) {
	cmd := NewExecuteSQLCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type ExecuteSQLCB func(result *ExecuteSQLResult, err error)

type AsyncExecuteSQLCommand struct {
	params *ExecuteSQLParams
	cb     ExecuteSQLCB
}

func NewAsyncExecuteSQLCommand(params *ExecuteSQLParams, cb ExecuteSQLCB) *AsyncExecuteSQLCommand {
	return &AsyncExecuteSQLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncExecuteSQLCommand) Name() string {
	return "Database.executeSQL"
}

func (cmd *AsyncExecuteSQLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ExecuteSQLCommand) Result() *ExecuteSQLResult {
	return &cmd.result
}

func (cmd *ExecuteSQLCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncExecuteSQLCommand) Done(data []byte, err error) {
	var result ExecuteSQLResult
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

type AddDatabaseEvent struct {
	Database *Database `json:"database"`
}

func OnAddDatabase(conn *hc.Conn, cb func(evt *AddDatabaseEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &AddDatabaseEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Database.addDatabase", sink)
}
