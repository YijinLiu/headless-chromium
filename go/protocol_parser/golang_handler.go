package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/yijinliu/algo-lib/go/src/logging"
)

type GolangProtocolHandler struct {
	outputDir  string
	handleExpr bool
	gofmt      string

	curVersion  string
	domains     []*ProtocolDomain
	nameCounts  map[string]int
	imports     map[string]string
	simpleTypes map[string]bool
}

func NewGolangProtocolHandler(outputDir string, handleExpr bool) *GolangProtocolHandler {
	gofmt, err := exec.LookPath("gofmt")
	if err != nil {
		logging.Vlog(0, "Failed to find gofmt binary. Will not run gofmt on generated go files.")
	}
	return &GolangProtocolHandler{
		outputDir:  outputDir,
		handleExpr: handleExpr,
		gofmt:      gofmt,
	}
}

func (h *GolangProtocolHandler) StartProtocol(version string) {
	h.curVersion = version
	h.domains = nil
	h.nameCounts = make(map[string]int)
	h.simpleTypes = make(map[string]bool)
}

func (h *GolangProtocolHandler) OnDomain(domain *ProtocolDomain) {
	h.domains = append(h.domains, domain)
	for _, tp := range domain.Types {
		name := toGolangType(tp.Id)
		h.nameCounts[name]++
		if tp.Type != "object" {
			h.simpleTypes[name] = true
		}
	}
	for _, cmd := range domain.Commands {
		name := toGolangType(cmd.Name)
		h.nameCounts[name]++
	}
	for _, evt := range domain.Events {
		name := toGolangType(evt.Name)
		h.nameCounts[name]++
	}
}

func (h *GolangProtocolHandler) EndProtocol() {
	for _, domain := range h.domains {
		h.processDomain(domain)
	}
}

func (h *GolangProtocolHandler) processDomain(domain *ProtocolDomain) {
	logging.Vlogf(2, "Processing domain %s ...", domain.Domain)
	if domain.Experimental && !h.handleExpr {
		logging.Vlogf(0, "Skip experimental domain '%s'.", domain.Domain)
		return
	}

	dir := filepath.Join(h.outputDir, "v"+h.curVersion)
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		logging.Fatal(err)
	}

	var buf bytes.Buffer
	h.imports = make(map[string]string)

	// Types.
	for _, tp := range domain.Types {
		if tp.Experimental && !h.handleExpr {
			logging.Vlogf(0, "\tSkip experimental type '%s'.", tp.Id)
		} else {
			logging.Vlogf(3, "\tProcessing type '%s' ...", tp.Id)
			h.onType(domain.Domain, tp, &buf)
		}
	}

	// Commands.
	for _, cmd := range domain.Commands {
		if cmd.Experimental && !h.handleExpr {
			logging.Vlogf(0, "\tSkip experimental command '%s'.", cmd.Name)
		} else {
			logging.Vlogf(3, "\tProcessing command '%s' ...", cmd.Name)
			h.onCommand(domain.Domain, cmd, &buf)
		}
	}

	// Events.
	if len(domain.Events) > 0 {
		for _, evt := range domain.Events {
			if evt.Experimental && !h.handleExpr {
				logging.Vlogf(0, "\tSkip experimental event '%s'.", evt.Name)
			} else {
				logging.Vlogf(3, "\tProcessing event '%s' ...", evt.Name)
				h.onEvent(domain.Domain, evt, &buf)
			}
		}
		h.imports["encoding/json"] = ""
		h.imports["github.com/yijinliu/algo-lib/go/src/logging"] = ""
	}
	h.writeGoFile(filepath.Join(dir, strings.ToLower(domain.Domain)+".go"), &buf)
}

func (h *GolangProtocolHandler) writeGoFile(file string, buf *bytes.Buffer) {
	f, err := os.Create(file)
	if err != nil {
		logging.Fatal(err)
	}
	if _, err := fmt.Fprintf(f, "package protocol\n\n"); err != nil {
		logging.Fatal(err)
	}
	if len(h.imports) > 0 {
		fmt.Fprintf(f, "import(\n")
		for path, name := range h.imports {
			if name != "" {
				fmt.Fprintf(f, "\t%s \"%s\"\n", name, path)
			} else {
				fmt.Fprintf(f, "\t\"%s\"\n", path)
			}
		}
		fmt.Fprintf(f, ")\n\n")
	}
	if _, err := buf.WriteTo(f); err != nil {
		logging.Fatal(err)
	}
	if h.gofmt == "" {
		return
	}
	cmd := exec.Command(h.gofmt, "-w", file)
	if err := cmd.Run(); err != nil {
		logging.Fatalf("Failed to run 'gofmt -w %s': %v", file, err)
	}
}

func experimentalTag(experimental bool) string {
	if experimental {
		return "// @experimental"
	}
	return ""
}

var descReplacer = strings.NewReplacer("<code>", "", "</code>", "")

func descriptionToGolangComment(desc string) string {
	if desc == "" {
		return ""
	}
	return "// " + descReplacer.Replace(desc)
}

func enumValueToGolangName(value string) string {
	start := true
	return strings.Map(func(ch rune) rune {
		if ch == '-' {
			if !start {
				start = true
				return -1
			}
			return 'N'
		}
		if start {
			start = false
			return unicode.ToUpper(ch)
		}
		return ch
	}, value)
}

func toGolangType(name string) string {
	return strings.Title(name)
}

func (h *GolangProtocolHandler) typeName(domain, name string) string {
	name = toGolangType(name)
	if h.nameCounts[name] > 1 {
		return domain + name
	}
	return name
}

var refReplacer = strings.NewReplacer(".", "")

func (h *GolangProtocolHandler) refToGolangType(domain, ref string) string {
	pos := strings.Index(ref, ".")
	var golangType string
	if pos == -1 {
		golangType = h.typeName(domain, ref)
	} else {
		golangType = h.typeName(ref[:pos], ref[pos+1:])
	}
	if h.simpleTypes[ref] {
		return golangType
	}
	return "*" + golangType
}

func (h *GolangProtocolHandler) simpleTypeToGolangType(domain string, st *SimpleType) string {
	switch st.Type {
	case "":
		if st.Ref == "" {
			logging.Fatalf("Illegal type '%v'.", st)
		}
		return h.refToGolangType(domain, st.Ref)
	case "number":
		return "float64"
	case "integer":
		return "int"
	case "any":
		h.imports["encoding/json"] = ""
		return "json.RawMessage"
	case "string":
		return "string"
	case "boolean":
		return "bool"
	case "object":
		return "map[string]string"
	}
	logging.Fatalf("Unknown type '%v'.", st)
	return ""
}

func (h *GolangProtocolHandler) unnamedTypeToGolangType(domain string, ut *UnnamedType) string {
	if ut.Type == "array" {
		return "[]" + h.simpleTypeToGolangType(domain, ut.Items)
	}
	return h.simpleTypeToGolangType(domain, &ut.SimpleType)
}

func (h *GolangProtocolHandler) onType(domain string, tp *DomainType, buf *bytes.Buffer) {
	name := h.typeName(domain, tp.Id)
	fmt.Fprintf(buf, "%s\n", descriptionToGolangComment(tp.Description))
	if tp.Experimental {
		fmt.Fprintln(buf, experimentalTag(tp.Experimental))
	}
	switch tp.Type {
	case "string":
		fmt.Fprintf(buf, "type %s string\n", name)
		// Define enum values.
		for _, val := range tp.Enum {
			fmt.Fprintf(buf, "const %s%s %s = \"%s\"\n", name,
				enumValueToGolangName(val), name, val)
		}
		buf.WriteRune('\n')
	case "object":
		fmt.Fprintf(buf, "type %s struct {\n", name)
		for _, prop := range tp.Properties {
			var omitEmpty string
			if prop.Optional {
				omitEmpty = ",omitempty"
			}
			fmt.Fprintf(buf, "\t%s %s `json:\"%s%s\"` %s\n", toGolangType(prop.Name),
				h.unnamedTypeToGolangType(domain, &prop.UnnamedType), prop.Name, omitEmpty,
				descriptionToGolangComment(prop.Description))
		}
		buf.WriteString("}\n\n")
	default:
		fmt.Fprintf(buf, "type %s %s\n\n", name,
			h.unnamedTypeToGolangType(domain, &tp.UnnamedType))
	}
}

func (h *GolangProtocolHandler) onCommand(domain string, cmd *DomainCommand, buf *bytes.Buffer) {
	h.imports["sync"] = ""
	h.imports["github.com/yijinliu/headless-chromium/go"] = "hc"
	name := h.typeName(domain, cmd.Name)

	// Params.
	var paramsField, paramsParam, paramsAssign, paramsValue, paramsName string
	if len(cmd.Parameters) > 0 {
		fmt.Fprintf(buf, "type %sParams struct {\n", name)
		for _, param := range cmd.Parameters {
			var omitEmpty string
			if param.Optional {
				omitEmpty = ",omitempty"
			}
			fmt.Fprintf(buf, "\t%s %s `json:\"%s%s\"` %s\n", toGolangType(param.Name),
				h.unnamedTypeToGolangType(domain, &param.UnnamedType), param.Name, omitEmpty,
				descriptionToGolangComment(param.Description))
		}
		buf.WriteString("}\n\n")
		paramsField = fmt.Sprintf("params *%sParams\n", name)
		paramsParam = fmt.Sprintf("params *%sParams, ", name)
		paramsAssign = "params: params,\n"
		paramsValue = "cmd.params"
		paramsName = "params"
	} else {
		paramsValue = "nil"
	}

	// Result.
	var resultField, resultParam, resultValue string
	if len(cmd.Returns) > 0 {
		fmt.Fprintf(buf, "type %sResult struct {\n", name)
		for _, ret := range cmd.Returns {
			fmt.Fprintf(buf, "\t%s %s `json:\"%s\"` %s\n", toGolangType(ret.Name),
				h.unnamedTypeToGolangType(domain, &ret.UnnamedType), ret.Name,
				descriptionToGolangComment(ret.Description))
		}
		buf.WriteString("}\n")
		resultField = fmt.Sprintf("result %sResult\n", name)
		resultParam = fmt.Sprintf("result *%sResult, ", name)
		resultValue = "&cmd.result, "
	}

	fmt.Fprintf(buf, `
%s
%s
type %sCommand struct {
	%s%swg sync.WaitGroup
	err error
}

func New%sCommand(%s) *%sCommand {
	return &%sCommand{
		%s
	}
}

func (cmd *%sCommand) Name() string {
	return "%s.%s"
}

func (cmd *%sCommand) Params() interface{} {
	return %s
}

func (cmd *%sCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func %s(%sconn *hc.Conn) (%serr error) {
	cmd := New%sCommand(%s)
	cmd.Run(conn)
	return %scmd.err
}

type %sCB func(%serr error)

%s
%s
type Async%sCommand struct {
	%scb %sCB
}

func NewAsync%sCommand(%scb %sCB) *Async%sCommand {
	return &Async%sCommand{
		%scb: cb,
	}
}

func (cmd *Async%sCommand) Name() string {
	return "%s.%s"
}

func (cmd *Async%sCommand) Params() interface{} {
	return %s
}
`,
		descriptionToGolangComment(cmd.Description), experimentalTag(cmd.Experimental), // comment
		name, paramsField, resultField, // struct
		name, paramsParam, name, name, paramsAssign, // constructor
		name, domain, cmd.Name, // method Name
		name, paramsValue, // method Params
		name,                                                          // method Run
		name, paramsParam, resultParam, name, paramsName, resultValue, // func Run
		name, resultParam, // CB
		descriptionToGolangComment(cmd.Description), experimentalTag(cmd.Experimental), // comment
		name, paramsField, name, // struct
		name, paramsParam, name, name, name, paramsAssign, // constructor
		name, domain, cmd.Name, // method Name
		name, paramsValue) // method Params

	if len(cmd.Returns) > 0 {
		fmt.Fprintf(buf, `
func (cmd *%sCommand) Result() *%sResult {
	return &cmd.result
}

func (cmd *%sCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *Async%sCommand) Done(data []byte, err error) {
	var result %sResult
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
`, name, name, name, name, name)
		h.imports["encoding/json"] = ""
		h.imports["github.com/yijinliu/algo-lib/go/src/logging"] = ""
	} else {
		fmt.Fprintf(buf, `
func (cmd *%sCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *Async%sCommand) Done(data []byte, err error) {
	cmd.cb(err)
}
`, name, name)
	}
}

func (h *GolangProtocolHandler) onEvent(domain string, evt *DomainEvent, buf *bytes.Buffer) {
	name := h.typeName(domain, evt.Name)

	// Params.
	fmt.Fprintf(buf, "%s\n%s\ntype %sEvent struct {\n", descriptionToGolangComment(evt.Description),
		experimentalTag(evt.Experimental), name)
	for _, param := range evt.Parameters {
		fmt.Fprintf(buf, "\t%s %s `json:\"%s\"` %s\n", toGolangType(param.Name),
			h.unnamedTypeToGolangType(domain, &param.UnnamedType), param.Name,
			descriptionToGolangComment(param.Description))
	}
	buf.WriteString("}\n\n")

	fmt.Fprintf(buf, `
func On%s(conn *hc.Conn, cb func(evt *%sEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &%sEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("%s.%s", sink)
}
`, name, name, name, domain, evt.Name)
}
