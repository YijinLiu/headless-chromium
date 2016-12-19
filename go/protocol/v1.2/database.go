package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

// Unique identifier of Database object.
type DatabaseId string

// Database object.
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

type DatabaseEnableCB func(err error)

// Enables database tracking, database events will now be delivered to the client.
type DatabaseEnableCommand struct {
	cb DatabaseEnableCB
}

func NewDatabaseEnableCommand(cb DatabaseEnableCB) *DatabaseEnableCommand {
	return &DatabaseEnableCommand{
		cb: cb,
	}
}

func (cmd *DatabaseEnableCommand) Name() string {
	return "Database.enable"
}

func (cmd *DatabaseEnableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type DatabaseDisableCB func(err error)

// Disables database tracking, prevents database events from being sent to the client.
type DatabaseDisableCommand struct {
	cb DatabaseDisableCB
}

func NewDatabaseDisableCommand(cb DatabaseDisableCB) *DatabaseDisableCommand {
	return &DatabaseDisableCommand{
		cb: cb,
	}
}

func (cmd *DatabaseDisableCommand) Name() string {
	return "Database.disable"
}

func (cmd *DatabaseDisableCommand) Params() interface{} {
	return nil
}

func (cmd *DatabaseDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetDatabaseTableNamesParams struct {
	DatabaseId DatabaseId `json:"databaseId"`
}

type GetDatabaseTableNamesResult struct {
	TableNames []string `json:"tableNames"`
}

type GetDatabaseTableNamesCB func(result *GetDatabaseTableNamesResult, err error)

type GetDatabaseTableNamesCommand struct {
	params *GetDatabaseTableNamesParams
	cb     GetDatabaseTableNamesCB
}

func NewGetDatabaseTableNamesCommand(params *GetDatabaseTableNamesParams, cb GetDatabaseTableNamesCB) *GetDatabaseTableNamesCommand {
	return &GetDatabaseTableNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetDatabaseTableNamesCommand) Name() string {
	return "Database.getDatabaseTableNames"
}

func (cmd *GetDatabaseTableNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetDatabaseTableNamesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetDatabaseTableNamesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ExecuteSQLParams struct {
	DatabaseId DatabaseId `json:"databaseId"`
	Query      string     `json:"query"`
}

type ExecuteSQLResult struct {
	ColumnNames []string `json:"columnNames"`
	Values      []string `json:"values"`
	SqlError    *Error   `json:"sqlError"`
}

type ExecuteSQLCB func(result *ExecuteSQLResult, err error)

type ExecuteSQLCommand struct {
	params *ExecuteSQLParams
	cb     ExecuteSQLCB
}

func NewExecuteSQLCommand(params *ExecuteSQLParams, cb ExecuteSQLCB) *ExecuteSQLCommand {
	return &ExecuteSQLCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ExecuteSQLCommand) Name() string {
	return "Database.executeSQL"
}

func (cmd *ExecuteSQLCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ExecuteSQLCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj ExecuteSQLResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type AddDatabaseEvent struct {
	Database *Database `json:"database"`
}

type AddDatabaseEventSink struct {
	events chan *AddDatabaseEvent
}

func NewAddDatabaseEventSink(bufSize int) *AddDatabaseEventSink {
	return &AddDatabaseEventSink{
		events: make(chan *AddDatabaseEvent, bufSize),
	}
}

func (s *AddDatabaseEventSink) Name() string {
	return "Database.addDatabase"
}

func (s *AddDatabaseEventSink) OnEvent(params []byte) {
	evt := &AddDatabaseEvent{}
	if err := json.Unmarshal(params, evt); err != nil {
		logging.Vlog(-1, err)
	} else {
		select {
		case s.events <- evt:
			// Do nothing.
		default:
			logging.Vlogf(0, "Dropped one event(%v).", evt)
		}
	}
}
