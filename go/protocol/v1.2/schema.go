package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// Description of the protocol domain.
type Domain struct {
	Name    string `json:"name"`    // Domain name.
	Version string `json:"version"` // Domain version.
}

type GetDomainsResult struct {
	Domains []*Domain `json:"domains"` // List of supported domains.
}

// Returns supported domains.

type GetDomainsCommand struct {
	result GetDomainsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetDomainsCommand() *GetDomainsCommand {
	return &GetDomainsCommand{}
}

func (cmd *GetDomainsCommand) Name() string {
	return "Schema.getDomains"
}

func (cmd *GetDomainsCommand) Params() interface{} {
	return nil
}

func (cmd *GetDomainsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetDomains(conn *hc.Conn) (result *GetDomainsResult, err error) {
	cmd := NewGetDomainsCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetDomainsCB func(result *GetDomainsResult, err error)

// Returns supported domains.

type AsyncGetDomainsCommand struct {
	cb GetDomainsCB
}

func NewAsyncGetDomainsCommand(cb GetDomainsCB) *AsyncGetDomainsCommand {
	return &AsyncGetDomainsCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetDomainsCommand) Name() string {
	return "Schema.getDomains"
}

func (cmd *AsyncGetDomainsCommand) Params() interface{} {
	return nil
}

func (cmd *GetDomainsCommand) Result() *GetDomainsResult {
	return &cmd.result
}

func (cmd *GetDomainsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetDomainsCommand) Done(data []byte, err error) {
	var result GetDomainsResult
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
