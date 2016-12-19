package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
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

type SecurityEnableCB func(err error)

// Enables tracking security state changes.
type SecurityEnableCommand struct {
	cb SecurityEnableCB
}

func NewSecurityEnableCommand(cb SecurityEnableCB) *SecurityEnableCommand {
	return &SecurityEnableCommand{
		cb: cb,
	}
}

func (cmd *SecurityEnableCommand) Name() string {
	return "Security.enable"
}

func (cmd *SecurityEnableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SecurityDisableCB func(err error)

// Disables tracking security state changes.
type SecurityDisableCommand struct {
	cb SecurityDisableCB
}

func NewSecurityDisableCommand(cb SecurityDisableCB) *SecurityDisableCommand {
	return &SecurityDisableCommand{
		cb: cb,
	}
}

func (cmd *SecurityDisableCommand) Name() string {
	return "Security.disable"
}

func (cmd *SecurityDisableCommand) Params() interface{} {
	return nil
}

func (cmd *SecurityDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type ShowCertificateViewerCB func(err error)

// Displays native dialog with the certificate details.
type ShowCertificateViewerCommand struct {
	cb ShowCertificateViewerCB
}

func NewShowCertificateViewerCommand(cb ShowCertificateViewerCB) *ShowCertificateViewerCommand {
	return &ShowCertificateViewerCommand{
		cb: cb,
	}
}

func (cmd *ShowCertificateViewerCommand) Name() string {
	return "Security.showCertificateViewer"
}

func (cmd *ShowCertificateViewerCommand) Params() interface{} {
	return nil
}

func (cmd *ShowCertificateViewerCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type SecurityStateChangedEvent struct {
	SecurityState         SecurityState               `json:"securityState"`         // Security state.
	Explanations          []*SecurityStateExplanation `json:"explanations"`          // List of explanations for the security state. If the overall security state is `insecure` or `warning`, at least one corresponding explanation should be included.
	InsecureContentStatus *InsecureContentStatus      `json:"insecureContentStatus"` // Information about insecure content on the page.
	SchemeIsCryptographic bool                        `json:"schemeIsCryptographic"` // True if the page was loaded over cryptographic transport such as HTTPS.
}

// The security state of the page changed.
type SecurityStateChangedEventSink struct {
	events chan *SecurityStateChangedEvent
}

func NewSecurityStateChangedEventSink(bufSize int) *SecurityStateChangedEventSink {
	return &SecurityStateChangedEventSink{
		events: make(chan *SecurityStateChangedEvent, bufSize),
	}
}

func (s *SecurityStateChangedEventSink) Name() string {
	return "Security.securityStateChanged"
}

func (s *SecurityStateChangedEventSink) OnEvent(params []byte) {
	evt := &SecurityStateChangedEvent{}
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
