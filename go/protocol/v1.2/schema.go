package protocol

import (
	"encoding/json"
)

// Description of the protocol domain.
type Domain struct {
	Name    string `json:"name"`    // Domain name.
	Version string `json:"version"` // Domain version.
}

type GetDomainsResult struct {
	Domains []*Domain `json:"domains"` // List of supported domains.
}

type GetDomainsCB func(result *GetDomainsResult, err error)

// Returns supported domains.
type GetDomainsCommand struct {
	cb GetDomainsCB
}

func NewGetDomainsCommand(cb GetDomainsCB) *GetDomainsCommand {
	return &GetDomainsCommand{
		cb: cb,
	}
}

func (cmd *GetDomainsCommand) Name() string {
	return "Schema.getDomains"
}

func (cmd *GetDomainsCommand) Params() interface{} {
	return nil
}

func (cmd *GetDomainsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetDomainsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}
