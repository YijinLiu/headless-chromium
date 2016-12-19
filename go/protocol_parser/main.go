package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yijinliu/algo-lib/go/src/logging"
)

var outputLangsFlag = flag.String("output-langs", "golang",
	"Languages separated by comma. Only golang supported for now.")

var golangOutputDirFlag = flag.String("golang-output-dir",
	"src/github.com/yijinliu/headless-chromium/go/protocol", "")
var golangHandleExperimentalFlag = flag.Bool("golang-handle-experimental", true, "")

func main() {
	flag.Parse()

	// Check command line.
	outputLangs := *outputLangsFlag
	if outputLangs == "" {
		logging.Fatal("Please specify --output-langs.")
	}
	phs := map[string]ProtocolHandler{}
	for _, lang := range strings.Split(outputLangs, ",") {
		switch lang {
		case "golang":
			phs[lang] =
				NewGolangProtocolHandler(*golangOutputDirFlag, *golangHandleExperimentalFlag)
		default:
			logging.Fatal("Unknown language: ", lang)
		}
	}

	// <protocol version => <domain name => domain> >
	protocolMap := make(map[string]map[string]*ProtocolDomain)

	// Parse protocol JSON definition files
	for _, pf := range flag.Args() {
		logging.Vlogf(1, "Processing '%s' ...", pf)
		var protocol Protocol
		if content, err := ioutil.ReadFile(pf); err != nil {
			logging.Fatal(err)
		} else if err := json.Unmarshal(content, &protocol); err != nil {
			logging.Fatal(err)
		}
		version := fmt.Sprintf("%s.%s", protocol.Version.Major, protocol.Version.Minor)
		domainMap := protocolMap[version]
		if domainMap == nil {
			domainMap = make(map[string]*ProtocolDomain)
			protocolMap[version] = domainMap
		}
		for _, domain := range protocol.Domains {
			if domainMap[domain.Domain] != nil {
				logging.Vlogf(0, "Domain '%s' is already defined!", domain.Domain)
			} else {
				domainMap[domain.Domain] = domain
			}
		}
	}

	// Process.
	for version, domainMap := range protocolMap {
		for lang, ph := range phs {
			logging.Vlogf(1, "Generating protocol v%s for %s ...", version, lang)
			ph.StartProtocol(version)
			for _, domain := range domainMap {
				ph.OnDomain(domain)
			}
			ph.EndProtocol()
		}
	}
}
