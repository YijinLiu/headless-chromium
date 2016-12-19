package main

type ProtocolVersion struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
}

type SimpleType struct {
	Type        string `json:"type"`
	Ref         string `json:"$ref"`
	Description string `json:"description"`
}

type UnnamedType struct {
	SimpleType
	Items    *SimpleType `json:"items"`
	Optional bool        `json:"optional"`
}

type NamedType struct {
	UnnamedType
	Name         string `json:"name"`
	Optional     bool   `json:"optional"`
	Experimental bool   `json:"experimental"`
}

type DomainType struct {
	UnnamedType
	Id           string       `json:"id"`
	Enum         []string     `json:"enum"`
	Properties   []*NamedType `json:"properties"`
	Experimental bool         `json:"experimental"`
	Exported     bool         `json:"exported"`
}

type DomainCommand struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Parameters   []*NamedType `json:"parameters"`
	Returns      []*NamedType `json:"returns"`
	Handlers     []string     `json:"handlers"`
	Redirect     string       `json:"redirect"`
	Experimental bool         `json:"experimental"`
	Async        bool         `json:"async"`
}

type DomainEvent struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Parameters   []*NamedType `json:"parameters"`
	Handlers     []string     `json:"handlers"`
	Experimental bool         `json:"experimental"`
}

type ProtocolDomain struct {
	Domain       string           `json:"domain"`
	Experimental bool             `json:"experimental"`
	Types        []*DomainType    `json:"types"`
	Commands     []*DomainCommand `json:"commands"`
	Events       []*DomainEvent   `json:"events"`
}

type Protocol struct {
	Version ProtocolVersion   `json:"version"`
	Domains []*ProtocolDomain `json:"domains"`
}

type ProtocolHandler interface {
	StartProtocol(version string)
	OnDomain(domain *ProtocolDomain)
	EndProtocol()
}
