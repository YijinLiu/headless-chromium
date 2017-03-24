package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

// An internal certificate ID value.
type CertificateId int

// The security level of a page or resource.
type SecurityState string

const SecurityStateUnknown SecurityState = "unknown"
const SecurityStateNeutral SecurityState = "neutral"
const SecurityStateInsecure SecurityState = "insecure"
const SecurityStateWarning SecurityState = "warning"
const SecurityStateSecure SecurityState = "secure"
const SecurityStateInfo SecurityState = "info"

// An explanation of an factor contributing to the security state.
type SecurityStateExplanation struct {
	SecurityState  SecurityState `json:"securityState"`  // Security state representing the severity of the factor being explained.
	Summary        string        `json:"summary"`        // Short phrase describing the type of factor.
	Description    string        `json:"description"`    // Full text explanation of the factor.
	HasCertificate bool          `json:"hasCertificate"` // True if the page has a certificate.
}

// Information about insecure content on the page.
type InsecureContentStatus struct {
	RanMixedContent                bool          `json:"ranMixedContent"`                // True if the page was loaded over HTTPS and ran mixed (HTTP) content such as scripts.
	DisplayedMixedContent          bool          `json:"displayedMixedContent"`          // True if the page was loaded over HTTPS and displayed mixed (HTTP) content such as images.
	RanContentWithCertErrors       bool          `json:"ranContentWithCertErrors"`       // True if the page was loaded over HTTPS without certificate errors, and ran content such as scripts that were loaded with certificate errors.
	DisplayedContentWithCertErrors bool          `json:"displayedContentWithCertErrors"` // True if the page was loaded over HTTPS without certificate errors, and displayed content such as images that were loaded with certificate errors.
	RanInsecureContentStyle        SecurityState `json:"ranInsecureContentStyle"`        // Security state representing a page that ran insecure content.
	DisplayedInsecureContentStyle  SecurityState `json:"displayedInsecureContentStyle"`  // Security state representing a page that displayed insecure content.
}

// Enables tracking security state changes.

type SecurityEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewSecurityEnableCommand() *SecurityEnableCommand {
	return &SecurityEnableCommand{}
}

func (cmd *SecurityEnableCommand) Name() string {
	return "Security.enable"
}

func (cmd *SecurityEnableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SecurityEnable(conn *hc.Conn) (err error) {
	cmd := NewSecurityEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type SecurityEnableCB func(err error)

// Enables tracking security state changes.

type AsyncSecurityEnableCommand struct {
	cb SecurityEnableCB
}

func NewAsyncSecurityEnableCommand(cb SecurityEnableCB) *AsyncSecurityEnableCommand {
	return &AsyncSecurityEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncSecurityEnableCommand) Name() string {
	return "Security.enable"
}

func (cmd *AsyncSecurityEnableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSecurityEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables tracking security state changes.

type SecurityDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewSecurityDisableCommand() *SecurityDisableCommand {
	return &SecurityDisableCommand{}
}

func (cmd *SecurityDisableCommand) Name() string {
	return "Security.disable"
}

func (cmd *SecurityDisableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SecurityDisable(conn *hc.Conn) (err error) {
	cmd := NewSecurityDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type SecurityDisableCB func(err error)

// Disables tracking security state changes.

type AsyncSecurityDisableCommand struct {
	cb SecurityDisableCB
}

func NewAsyncSecurityDisableCommand(cb SecurityDisableCB) *AsyncSecurityDisableCommand {
	return &AsyncSecurityDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncSecurityDisableCommand) Name() string {
	return "Security.disable"
}

func (cmd *AsyncSecurityDisableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSecurityDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Displays native dialog with the certificate details.

type ShowCertificateViewerCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewShowCertificateViewerCommand() *ShowCertificateViewerCommand {
	return &ShowCertificateViewerCommand{}
}

func (cmd *ShowCertificateViewerCommand) Name() string {
	return "Security.showCertificateViewer"
}

func (cmd *ShowCertificateViewerCommand) Params() interface{} {
	return nil
}

func (cmd *ShowCertificateViewerCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ShowCertificateViewer(conn *hc.Conn) (err error) {
	cmd := NewShowCertificateViewerCommand()
	cmd.Run(conn)
	return cmd.err
}

type ShowCertificateViewerCB func(err error)

// Displays native dialog with the certificate details.

type AsyncShowCertificateViewerCommand struct {
	cb ShowCertificateViewerCB
}

func NewAsyncShowCertificateViewerCommand(cb ShowCertificateViewerCB) *AsyncShowCertificateViewerCommand {
	return &AsyncShowCertificateViewerCommand{
		cb: cb,
	}
}

func (cmd *AsyncShowCertificateViewerCommand) Name() string {
	return "Security.showCertificateViewer"
}

func (cmd *AsyncShowCertificateViewerCommand) Params() interface{} {
	return nil
}

func (cmd *ShowCertificateViewerCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncShowCertificateViewerCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// The security state of the page changed.

type SecurityStateChangedEvent struct {
	SecurityState         SecurityState               `json:"securityState"`         // Security state.
	SchemeIsCryptographic bool                        `json:"schemeIsCryptographic"` // True if the page was loaded over cryptographic transport such as HTTPS.
	Explanations          []*SecurityStateExplanation `json:"explanations"`          // List of explanations for the security state. If the overall security state is `insecure` or `warning`, at least one corresponding explanation should be included.
	InsecureContentStatus *InsecureContentStatus      `json:"insecureContentStatus"` // Information about insecure content on the page.
	Summary               string                      `json:"summary"`               // Overrides user-visible description of the state.
}

func OnSecurityStateChanged(conn *hc.Conn, cb func(evt *SecurityStateChangedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &SecurityStateChangedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("Security.securityStateChanged", sink)
}
