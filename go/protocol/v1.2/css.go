package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
	hc "github.com/yijinliu/headless-chromium/go"
	"sync"
)

type StyleSheetId string

// Stylesheet type: "injected" for stylesheets injected via extension, "user-agent" for user-agent stylesheets, "inspector" for stylesheets created by the inspector (i.e. those holding the "via inspector" rules), "regular" for regular stylesheets.
type StyleSheetOrigin string

const StyleSheetOriginInjected StyleSheetOrigin = "injected"
const StyleSheetOriginUserAgent StyleSheetOrigin = "user-agent"
const StyleSheetOriginInspector StyleSheetOrigin = "inspector"
const StyleSheetOriginRegular StyleSheetOrigin = "regular"

// CSS rule collection for a single pseudo style.
type PseudoElementMatches struct {
	PseudoType *PseudoType  `json:"pseudoType"` // Pseudo element type.
	Matches    []*RuleMatch `json:"matches"`    // Matches of CSS rules applicable to the pseudo style.
}

// Inherited CSS rule collection from ancestor node.
type InheritedStyleEntry struct {
	InlineStyle     *CSSStyle    `json:"inlineStyle,omitempty"` // The ancestor node's inline style, if any, in the style inheritance chain.
	MatchedCSSRules []*RuleMatch `json:"matchedCSSRules"`       // Matches of CSS rules matching the ancestor node in the style inheritance chain.
}

// Match data for a CSS rule.
type RuleMatch struct {
	Rule              *CSSRule `json:"rule"`              // CSS rule in the match.
	MatchingSelectors []int    `json:"matchingSelectors"` // Matching selector indices in the rule's selectorList selectors (0-based).
}

// Data for a simple selector (these are delimited by commas in a selector list).
type Value struct {
	Text  string       `json:"text"`            // Value text.
	Range *SourceRange `json:"range,omitempty"` // Value range in the underlying resource (if available).
}

// Selector list data.
type SelectorList struct {
	Selectors []*Value `json:"selectors"` // Selectors in the list.
	Text      string   `json:"text"`      // Rule selector text.
}

// CSS stylesheet metainformation.
type CSSStyleSheetHeader struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId"`           // The stylesheet identifier.
	FrameId      *FrameId         `json:"frameId"`                // Owner frame identifier.
	SourceURL    string           `json:"sourceURL"`              // Stylesheet resource URL.
	SourceMapURL string           `json:"sourceMapURL,omitempty"` // URL of source map associated with the stylesheet (if any).
	Origin       StyleSheetOrigin `json:"origin"`                 // Stylesheet origin.
	Title        string           `json:"title"`                  // Stylesheet title.
	OwnerNode    *BackendNodeId   `json:"ownerNode,omitempty"`    // The backend id for the owner node of the stylesheet.
	Disabled     bool             `json:"disabled"`               // Denotes whether the stylesheet is disabled.
	HasSourceURL bool             `json:"hasSourceURL,omitempty"` // Whether the sourceURL field value comes from the sourceURL comment.
	IsInline     bool             `json:"isInline"`               // Whether this stylesheet is created for STYLE tag by parser. This flag is not set for document.written STYLE tags.
	StartLine    float64          `json:"startLine"`              // Line offset of the stylesheet within the resource (zero based).
	StartColumn  float64          `json:"startColumn"`            // Column offset of the stylesheet within the resource (zero based).
}

// CSS rule representation.
type CSSRule struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId,omitempty"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	SelectorList *SelectorList    `json:"selectorList"`           // Rule selector data.
	Origin       StyleSheetOrigin `json:"origin"`                 // Parent stylesheet's origin.
	Style        *CSSStyle        `json:"style"`                  // Associated style declaration.
	Media        []*CSSMedia      `json:"media,omitempty"`        // Media list array (for rules involving media queries). The array enumerates media queries starting with the innermost one, going outwards.
}

// CSS rule usage information.
// @experimental
type RuleUsage struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	Range        *SourceRange `json:"range"`        // Style declaration range in the enclosing stylesheet (if available).
	Used         bool         `json:"used"`         // Indicates whether the rule was actually used by some element in the page.
}

// Text range within a resource. All numbers are zero-based.
type SourceRange struct {
	StartLine   int `json:"startLine"`   // Start line of range.
	StartColumn int `json:"startColumn"` // Start column of range (inclusive).
	EndLine     int `json:"endLine"`     // End line of range
	EndColumn   int `json:"endColumn"`   // End column of range (exclusive).
}

type ShorthandEntry struct {
	Name      string `json:"name"`                // Shorthand name.
	Value     string `json:"value"`               // Shorthand value.
	Important bool   `json:"important,omitempty"` // Whether the property has "!important" annotation (implies false if absent).
}

type CSSComputedStyleProperty struct {
	Name  string `json:"name"`  // Computed style property name.
	Value string `json:"value"` // Computed style property value.
}

// CSS style representation.
type CSSStyle struct {
	StyleSheetId     StyleSheetId      `json:"styleSheetId,omitempty"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	CssProperties    []*CSSProperty    `json:"cssProperties"`          // CSS properties in the style.
	ShorthandEntries []*ShorthandEntry `json:"shorthandEntries"`       // Computed values for all shorthands found in the style.
	CssText          string            `json:"cssText,omitempty"`      // Style declaration text (if available).
	Range            *SourceRange      `json:"range,omitempty"`        // Style declaration range in the enclosing stylesheet (if available).
}

// CSS property declaration data.
type CSSProperty struct {
	Name      string       `json:"name"`                // The property name.
	Value     string       `json:"value"`               // The property value.
	Important bool         `json:"important,omitempty"` // Whether the property has "!important" annotation (implies false if absent).
	Implicit  bool         `json:"implicit,omitempty"`  // Whether the property is implicit (implies false if absent).
	Text      string       `json:"text,omitempty"`      // The full property text as specified in the style.
	ParsedOk  bool         `json:"parsedOk,omitempty"`  // Whether the property is understood by the browser (implies true if absent).
	Disabled  bool         `json:"disabled,omitempty"`  // Whether the property is disabled by the user (present for source-based properties only).
	Range     *SourceRange `json:"range,omitempty"`     // The entire property range in the enclosing style declaration (if available).
}

// CSS media rule descriptor.
type CSSMedia struct {
	Text         string        `json:"text"`                   // Media query text.
	Source       string        `json:"source"`                 // Source of the media query: "mediaRule" if specified by a @media rule, "importRule" if specified by an @import rule, "linkedSheet" if specified by a "media" attribute in a linked stylesheet's LINK tag, "inlineSheet" if specified by a "media" attribute in an inline stylesheet's STYLE tag.
	SourceURL    string        `json:"sourceURL,omitempty"`    // URL of the document containing the media query description.
	Range        *SourceRange  `json:"range,omitempty"`        // The associated rule (@media or @import) header range in the enclosing stylesheet (if available).
	StyleSheetId StyleSheetId  `json:"styleSheetId,omitempty"` // Identifier of the stylesheet containing this object (if exists).
	MediaList    []*MediaQuery `json:"mediaList,omitempty"`    // Array of media queries.
}

// Media query descriptor.
// @experimental
type MediaQuery struct {
	Expressions []*MediaQueryExpression `json:"expressions"` // Array of media query expressions.
	Active      bool                    `json:"active"`      // Whether the media query condition is satisfied.
}

// Media query expression descriptor.
// @experimental
type MediaQueryExpression struct {
	Value          float64      `json:"value"`                    // Media query expression value.
	Unit           string       `json:"unit"`                     // Media query expression units.
	Feature        string       `json:"feature"`                  // Media query expression feature.
	ValueRange     *SourceRange `json:"valueRange,omitempty"`     // The associated range of the value text in the enclosing stylesheet (if available).
	ComputedLength float64      `json:"computedLength,omitempty"` // Computed length of media query expression (if applicable).
}

// Information about amount of glyphs that were rendered with given font.
// @experimental
type PlatformFontUsage struct {
	FamilyName   string  `json:"familyName"`   // Font's family name reported by platform.
	IsCustomFont bool    `json:"isCustomFont"` // Indicates if the font was downloaded or resolved locally.
	GlyphCount   float64 `json:"glyphCount"`   // Amount of glyphs that were rendered with this font.
}

// CSS keyframes rule representation.
type CSSKeyframesRule struct {
	AnimationName *Value             `json:"animationName"` // Animation name.
	Keyframes     []*CSSKeyframeRule `json:"keyframes"`     // List of keyframes.
}

// CSS keyframe rule representation.
type CSSKeyframeRule struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId,omitempty"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	Origin       StyleSheetOrigin `json:"origin"`                 // Parent stylesheet's origin.
	KeyText      *Value           `json:"keyText"`                // Associated key text.
	Style        *CSSStyle        `json:"style"`                  // Associated style declaration.
}

// A descriptor of operation to mutate style declaration text.
type StyleDeclarationEdit struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // The css style sheet identifier.
	Range        *SourceRange `json:"range"`        // The range of the style text in the enclosing stylesheet.
	Text         string       `json:"text"`         // New style text.
}

// Details of post layout rendered text positions. The exact layout should not be regarded as stable and may change between versions.
// @experimental
type InlineTextBox struct {
	BoundingBox         *Rect `json:"boundingBox"`         // The absolute position bounding box.
	StartCharacterIndex int   `json:"startCharacterIndex"` // The starting index in characters, for this post layout textbox substring.
	NumCharacters       int   `json:"numCharacters"`       // The number of characters in this post layout textbox substring.
}

// Details of an element in the DOM tree with a LayoutObject.
// @experimental
type LayoutTreeNode struct {
	NodeId          *NodeId          `json:"nodeId"`                    // The id of the related DOM node matching one from DOM.GetDocument.
	BoundingBox     *Rect            `json:"boundingBox"`               // The absolute position bounding box.
	LayoutText      string           `json:"layoutText,omitempty"`      // Contents of the LayoutText if any
	InlineTextNodes []*InlineTextBox `json:"inlineTextNodes,omitempty"` // The post layout inline text nodes, if any.
	StyleIndex      int              `json:"styleIndex,omitempty"`      // Index into the computedStyles array returned by getLayoutTreeAndStyles.
}

// A subset of the full ComputedStyle as defined by the request whitelist.
// @experimental
type ComputedStyle struct {
	Properties []*CSSComputedStyleProperty `json:"properties"`
}

// Enables the CSS agent for the given page. Clients should not assume that the CSS agent has been enabled until the result of this command is received.

type CSSEnableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewCSSEnableCommand() *CSSEnableCommand {
	return &CSSEnableCommand{}
}

func (cmd *CSSEnableCommand) Name() string {
	return "CSS.enable"
}

func (cmd *CSSEnableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSEnableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CSSEnable(conn *hc.Conn) (err error) {
	cmd := NewCSSEnableCommand()
	cmd.Run(conn)
	return cmd.err
}

type CSSEnableCB func(err error)

// Enables the CSS agent for the given page. Clients should not assume that the CSS agent has been enabled until the result of this command is received.

type AsyncCSSEnableCommand struct {
	cb CSSEnableCB
}

func NewAsyncCSSEnableCommand(cb CSSEnableCB) *AsyncCSSEnableCommand {
	return &AsyncCSSEnableCommand{
		cb: cb,
	}
}

func (cmd *AsyncCSSEnableCommand) Name() string {
	return "CSS.enable"
}

func (cmd *AsyncCSSEnableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSEnableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCSSEnableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

// Disables the CSS agent for the given page.

type CSSDisableCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewCSSDisableCommand() *CSSDisableCommand {
	return &CSSDisableCommand{}
}

func (cmd *CSSDisableCommand) Name() string {
	return "CSS.disable"
}

func (cmd *CSSDisableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSDisableCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CSSDisable(conn *hc.Conn) (err error) {
	cmd := NewCSSDisableCommand()
	cmd.Run(conn)
	return cmd.err
}

type CSSDisableCB func(err error)

// Disables the CSS agent for the given page.

type AsyncCSSDisableCommand struct {
	cb CSSDisableCB
}

func NewAsyncCSSDisableCommand(cb CSSDisableCB) *AsyncCSSDisableCommand {
	return &AsyncCSSDisableCommand{
		cb: cb,
	}
}

func (cmd *AsyncCSSDisableCommand) Name() string {
	return "CSS.disable"
}

func (cmd *AsyncCSSDisableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSDisableCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCSSDisableCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetMatchedStylesForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetMatchedStylesForNodeResult struct {
	InlineStyle       *CSSStyle               `json:"inlineStyle"`       // Inline style for the specified DOM node.
	AttributesStyle   *CSSStyle               `json:"attributesStyle"`   // Attribute-defined element style (e.g. resulting from "width=20 height=100%").
	MatchedCSSRules   []*RuleMatch            `json:"matchedCSSRules"`   // CSS rules matching this node, from all applicable stylesheets.
	PseudoElements    []*PseudoElementMatches `json:"pseudoElements"`    // Pseudo style matches for this node.
	Inherited         []*InheritedStyleEntry  `json:"inherited"`         // A chain of inherited styles (from the immediate node parent up to the DOM tree root).
	CssKeyframesRules []*CSSKeyframesRule     `json:"cssKeyframesRules"` // A list of CSS keyframed animations matching this node.
}

// Returns requested styles for a DOM node identified by nodeId.

type GetMatchedStylesForNodeCommand struct {
	params *GetMatchedStylesForNodeParams
	result GetMatchedStylesForNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetMatchedStylesForNodeCommand(params *GetMatchedStylesForNodeParams) *GetMatchedStylesForNodeCommand {
	return &GetMatchedStylesForNodeCommand{
		params: params,
	}
}

func (cmd *GetMatchedStylesForNodeCommand) Name() string {
	return "CSS.getMatchedStylesForNode"
}

func (cmd *GetMatchedStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetMatchedStylesForNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetMatchedStylesForNode(params *GetMatchedStylesForNodeParams, conn *hc.Conn) (result *GetMatchedStylesForNodeResult, err error) {
	cmd := NewGetMatchedStylesForNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetMatchedStylesForNodeCB func(result *GetMatchedStylesForNodeResult, err error)

// Returns requested styles for a DOM node identified by nodeId.

type AsyncGetMatchedStylesForNodeCommand struct {
	params *GetMatchedStylesForNodeParams
	cb     GetMatchedStylesForNodeCB
}

func NewAsyncGetMatchedStylesForNodeCommand(params *GetMatchedStylesForNodeParams, cb GetMatchedStylesForNodeCB) *AsyncGetMatchedStylesForNodeCommand {
	return &AsyncGetMatchedStylesForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetMatchedStylesForNodeCommand) Name() string {
	return "CSS.getMatchedStylesForNode"
}

func (cmd *AsyncGetMatchedStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetMatchedStylesForNodeCommand) Result() *GetMatchedStylesForNodeResult {
	return &cmd.result
}

func (cmd *GetMatchedStylesForNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetMatchedStylesForNodeCommand) Done(data []byte, err error) {
	var result GetMatchedStylesForNodeResult
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

type GetInlineStylesForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetInlineStylesForNodeResult struct {
	InlineStyle     *CSSStyle `json:"inlineStyle"`     // Inline style for the specified DOM node.
	AttributesStyle *CSSStyle `json:"attributesStyle"` // Attribute-defined element style (e.g. resulting from "width=20 height=100%").
}

// Returns the styles defined inline (explicitly in the "style" attribute and implicitly, using DOM attributes) for a DOM node identified by nodeId.

type GetInlineStylesForNodeCommand struct {
	params *GetInlineStylesForNodeParams
	result GetInlineStylesForNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetInlineStylesForNodeCommand(params *GetInlineStylesForNodeParams) *GetInlineStylesForNodeCommand {
	return &GetInlineStylesForNodeCommand{
		params: params,
	}
}

func (cmd *GetInlineStylesForNodeCommand) Name() string {
	return "CSS.getInlineStylesForNode"
}

func (cmd *GetInlineStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetInlineStylesForNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetInlineStylesForNode(params *GetInlineStylesForNodeParams, conn *hc.Conn) (result *GetInlineStylesForNodeResult, err error) {
	cmd := NewGetInlineStylesForNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetInlineStylesForNodeCB func(result *GetInlineStylesForNodeResult, err error)

// Returns the styles defined inline (explicitly in the "style" attribute and implicitly, using DOM attributes) for a DOM node identified by nodeId.

type AsyncGetInlineStylesForNodeCommand struct {
	params *GetInlineStylesForNodeParams
	cb     GetInlineStylesForNodeCB
}

func NewAsyncGetInlineStylesForNodeCommand(params *GetInlineStylesForNodeParams, cb GetInlineStylesForNodeCB) *AsyncGetInlineStylesForNodeCommand {
	return &AsyncGetInlineStylesForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetInlineStylesForNodeCommand) Name() string {
	return "CSS.getInlineStylesForNode"
}

func (cmd *AsyncGetInlineStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetInlineStylesForNodeCommand) Result() *GetInlineStylesForNodeResult {
	return &cmd.result
}

func (cmd *GetInlineStylesForNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetInlineStylesForNodeCommand) Done(data []byte, err error) {
	var result GetInlineStylesForNodeResult
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

type GetComputedStyleForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetComputedStyleForNodeResult struct {
	ComputedStyle []*CSSComputedStyleProperty `json:"computedStyle"` // Computed style for the specified DOM node.
}

// Returns the computed style for a DOM node identified by nodeId.

type GetComputedStyleForNodeCommand struct {
	params *GetComputedStyleForNodeParams
	result GetComputedStyleForNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetComputedStyleForNodeCommand(params *GetComputedStyleForNodeParams) *GetComputedStyleForNodeCommand {
	return &GetComputedStyleForNodeCommand{
		params: params,
	}
}

func (cmd *GetComputedStyleForNodeCommand) Name() string {
	return "CSS.getComputedStyleForNode"
}

func (cmd *GetComputedStyleForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetComputedStyleForNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetComputedStyleForNode(params *GetComputedStyleForNodeParams, conn *hc.Conn) (result *GetComputedStyleForNodeResult, err error) {
	cmd := NewGetComputedStyleForNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetComputedStyleForNodeCB func(result *GetComputedStyleForNodeResult, err error)

// Returns the computed style for a DOM node identified by nodeId.

type AsyncGetComputedStyleForNodeCommand struct {
	params *GetComputedStyleForNodeParams
	cb     GetComputedStyleForNodeCB
}

func NewAsyncGetComputedStyleForNodeCommand(params *GetComputedStyleForNodeParams, cb GetComputedStyleForNodeCB) *AsyncGetComputedStyleForNodeCommand {
	return &AsyncGetComputedStyleForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetComputedStyleForNodeCommand) Name() string {
	return "CSS.getComputedStyleForNode"
}

func (cmd *AsyncGetComputedStyleForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetComputedStyleForNodeCommand) Result() *GetComputedStyleForNodeResult {
	return &cmd.result
}

func (cmd *GetComputedStyleForNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetComputedStyleForNodeCommand) Done(data []byte, err error) {
	var result GetComputedStyleForNodeResult
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

type GetPlatformFontsForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetPlatformFontsForNodeResult struct {
	Fonts []*PlatformFontUsage `json:"fonts"` // Usage statistics for every employed platform font.
}

// Requests information about platform fonts which we used to render child TextNodes in the given node.
// @experimental
type GetPlatformFontsForNodeCommand struct {
	params *GetPlatformFontsForNodeParams
	result GetPlatformFontsForNodeResult
	wg     sync.WaitGroup
	err    error
}

func NewGetPlatformFontsForNodeCommand(params *GetPlatformFontsForNodeParams) *GetPlatformFontsForNodeCommand {
	return &GetPlatformFontsForNodeCommand{
		params: params,
	}
}

func (cmd *GetPlatformFontsForNodeCommand) Name() string {
	return "CSS.getPlatformFontsForNode"
}

func (cmd *GetPlatformFontsForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPlatformFontsForNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetPlatformFontsForNode(params *GetPlatformFontsForNodeParams, conn *hc.Conn) (result *GetPlatformFontsForNodeResult, err error) {
	cmd := NewGetPlatformFontsForNodeCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetPlatformFontsForNodeCB func(result *GetPlatformFontsForNodeResult, err error)

// Requests information about platform fonts which we used to render child TextNodes in the given node.
// @experimental
type AsyncGetPlatformFontsForNodeCommand struct {
	params *GetPlatformFontsForNodeParams
	cb     GetPlatformFontsForNodeCB
}

func NewAsyncGetPlatformFontsForNodeCommand(params *GetPlatformFontsForNodeParams, cb GetPlatformFontsForNodeCB) *AsyncGetPlatformFontsForNodeCommand {
	return &AsyncGetPlatformFontsForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetPlatformFontsForNodeCommand) Name() string {
	return "CSS.getPlatformFontsForNode"
}

func (cmd *AsyncGetPlatformFontsForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPlatformFontsForNodeCommand) Result() *GetPlatformFontsForNodeResult {
	return &cmd.result
}

func (cmd *GetPlatformFontsForNodeCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetPlatformFontsForNodeCommand) Done(data []byte, err error) {
	var result GetPlatformFontsForNodeResult
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

type GetStyleSheetTextParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

type GetStyleSheetTextResult struct {
	Text string `json:"text"` // The stylesheet text.
}

// Returns the current textual content and the URL for a stylesheet.

type GetStyleSheetTextCommand struct {
	params *GetStyleSheetTextParams
	result GetStyleSheetTextResult
	wg     sync.WaitGroup
	err    error
}

func NewGetStyleSheetTextCommand(params *GetStyleSheetTextParams) *GetStyleSheetTextCommand {
	return &GetStyleSheetTextCommand{
		params: params,
	}
}

func (cmd *GetStyleSheetTextCommand) Name() string {
	return "CSS.getStyleSheetText"
}

func (cmd *GetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetStyleSheetTextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetStyleSheetText(params *GetStyleSheetTextParams, conn *hc.Conn) (result *GetStyleSheetTextResult, err error) {
	cmd := NewGetStyleSheetTextCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetStyleSheetTextCB func(result *GetStyleSheetTextResult, err error)

// Returns the current textual content and the URL for a stylesheet.

type AsyncGetStyleSheetTextCommand struct {
	params *GetStyleSheetTextParams
	cb     GetStyleSheetTextCB
}

func NewAsyncGetStyleSheetTextCommand(params *GetStyleSheetTextParams, cb GetStyleSheetTextCB) *AsyncGetStyleSheetTextCommand {
	return &AsyncGetStyleSheetTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetStyleSheetTextCommand) Name() string {
	return "CSS.getStyleSheetText"
}

func (cmd *AsyncGetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetStyleSheetTextCommand) Result() *GetStyleSheetTextResult {
	return &cmd.result
}

func (cmd *GetStyleSheetTextCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetStyleSheetTextCommand) Done(data []byte, err error) {
	var result GetStyleSheetTextResult
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

type CollectClassNamesParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

type CollectClassNamesResult struct {
	ClassNames []string `json:"classNames"` // Class name list.
}

// Returns all class names from specified stylesheet.
// @experimental
type CollectClassNamesCommand struct {
	params *CollectClassNamesParams
	result CollectClassNamesResult
	wg     sync.WaitGroup
	err    error
}

func NewCollectClassNamesCommand(params *CollectClassNamesParams) *CollectClassNamesCommand {
	return &CollectClassNamesCommand{
		params: params,
	}
}

func (cmd *CollectClassNamesCommand) Name() string {
	return "CSS.collectClassNames"
}

func (cmd *CollectClassNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CollectClassNames(params *CollectClassNamesParams, conn *hc.Conn) (result *CollectClassNamesResult, err error) {
	cmd := NewCollectClassNamesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CollectClassNamesCB func(result *CollectClassNamesResult, err error)

// Returns all class names from specified stylesheet.
// @experimental
type AsyncCollectClassNamesCommand struct {
	params *CollectClassNamesParams
	cb     CollectClassNamesCB
}

func NewAsyncCollectClassNamesCommand(params *CollectClassNamesParams, cb CollectClassNamesCB) *AsyncCollectClassNamesCommand {
	return &AsyncCollectClassNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCollectClassNamesCommand) Name() string {
	return "CSS.collectClassNames"
}

func (cmd *AsyncCollectClassNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesCommand) Result() *CollectClassNamesResult {
	return &cmd.result
}

func (cmd *CollectClassNamesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCollectClassNamesCommand) Done(data []byte, err error) {
	var result CollectClassNamesResult
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

type SetStyleSheetTextParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
	Text         string       `json:"text"`
}

type SetStyleSheetTextResult struct {
	SourceMapURL string `json:"sourceMapURL"` // URL of source map associated with script (if any).
}

// Sets the new stylesheet text.

type SetStyleSheetTextCommand struct {
	params *SetStyleSheetTextParams
	result SetStyleSheetTextResult
	wg     sync.WaitGroup
	err    error
}

func NewSetStyleSheetTextCommand(params *SetStyleSheetTextParams) *SetStyleSheetTextCommand {
	return &SetStyleSheetTextCommand{
		params: params,
	}
}

func (cmd *SetStyleSheetTextCommand) Name() string {
	return "CSS.setStyleSheetText"
}

func (cmd *SetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleSheetTextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetStyleSheetText(params *SetStyleSheetTextParams, conn *hc.Conn) (result *SetStyleSheetTextResult, err error) {
	cmd := NewSetStyleSheetTextCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetStyleSheetTextCB func(result *SetStyleSheetTextResult, err error)

// Sets the new stylesheet text.

type AsyncSetStyleSheetTextCommand struct {
	params *SetStyleSheetTextParams
	cb     SetStyleSheetTextCB
}

func NewAsyncSetStyleSheetTextCommand(params *SetStyleSheetTextParams, cb SetStyleSheetTextCB) *AsyncSetStyleSheetTextCommand {
	return &AsyncSetStyleSheetTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetStyleSheetTextCommand) Name() string {
	return "CSS.setStyleSheetText"
}

func (cmd *AsyncSetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleSheetTextCommand) Result() *SetStyleSheetTextResult {
	return &cmd.result
}

func (cmd *SetStyleSheetTextCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetStyleSheetTextCommand) Done(data []byte, err error) {
	var result SetStyleSheetTextResult
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

type SetRuleSelectorParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
	Range        *SourceRange `json:"range"`
	Selector     string       `json:"selector"`
}

type SetRuleSelectorResult struct {
	SelectorList *SelectorList `json:"selectorList"` // The resulting selector list after modification.
}

// Modifies the rule selector.

type SetRuleSelectorCommand struct {
	params *SetRuleSelectorParams
	result SetRuleSelectorResult
	wg     sync.WaitGroup
	err    error
}

func NewSetRuleSelectorCommand(params *SetRuleSelectorParams) *SetRuleSelectorCommand {
	return &SetRuleSelectorCommand{
		params: params,
	}
}

func (cmd *SetRuleSelectorCommand) Name() string {
	return "CSS.setRuleSelector"
}

func (cmd *SetRuleSelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRuleSelectorCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetRuleSelector(params *SetRuleSelectorParams, conn *hc.Conn) (result *SetRuleSelectorResult, err error) {
	cmd := NewSetRuleSelectorCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetRuleSelectorCB func(result *SetRuleSelectorResult, err error)

// Modifies the rule selector.

type AsyncSetRuleSelectorCommand struct {
	params *SetRuleSelectorParams
	cb     SetRuleSelectorCB
}

func NewAsyncSetRuleSelectorCommand(params *SetRuleSelectorParams, cb SetRuleSelectorCB) *AsyncSetRuleSelectorCommand {
	return &AsyncSetRuleSelectorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetRuleSelectorCommand) Name() string {
	return "CSS.setRuleSelector"
}

func (cmd *AsyncSetRuleSelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRuleSelectorCommand) Result() *SetRuleSelectorResult {
	return &cmd.result
}

func (cmd *SetRuleSelectorCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetRuleSelectorCommand) Done(data []byte, err error) {
	var result SetRuleSelectorResult
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

type SetKeyframeKeyParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
	Range        *SourceRange `json:"range"`
	KeyText      string       `json:"keyText"`
}

type SetKeyframeKeyResult struct {
	KeyText *Value `json:"keyText"` // The resulting key text after modification.
}

// Modifies the keyframe rule key text.

type SetKeyframeKeyCommand struct {
	params *SetKeyframeKeyParams
	result SetKeyframeKeyResult
	wg     sync.WaitGroup
	err    error
}

func NewSetKeyframeKeyCommand(params *SetKeyframeKeyParams) *SetKeyframeKeyCommand {
	return &SetKeyframeKeyCommand{
		params: params,
	}
}

func (cmd *SetKeyframeKeyCommand) Name() string {
	return "CSS.setKeyframeKey"
}

func (cmd *SetKeyframeKeyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetKeyframeKeyCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetKeyframeKey(params *SetKeyframeKeyParams, conn *hc.Conn) (result *SetKeyframeKeyResult, err error) {
	cmd := NewSetKeyframeKeyCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetKeyframeKeyCB func(result *SetKeyframeKeyResult, err error)

// Modifies the keyframe rule key text.

type AsyncSetKeyframeKeyCommand struct {
	params *SetKeyframeKeyParams
	cb     SetKeyframeKeyCB
}

func NewAsyncSetKeyframeKeyCommand(params *SetKeyframeKeyParams, cb SetKeyframeKeyCB) *AsyncSetKeyframeKeyCommand {
	return &AsyncSetKeyframeKeyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetKeyframeKeyCommand) Name() string {
	return "CSS.setKeyframeKey"
}

func (cmd *AsyncSetKeyframeKeyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetKeyframeKeyCommand) Result() *SetKeyframeKeyResult {
	return &cmd.result
}

func (cmd *SetKeyframeKeyCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetKeyframeKeyCommand) Done(data []byte, err error) {
	var result SetKeyframeKeyResult
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

type SetStyleTextsParams struct {
	Edits []*StyleDeclarationEdit `json:"edits"`
}

type SetStyleTextsResult struct {
	Styles []*CSSStyle `json:"styles"` // The resulting styles after modification.
}

// Applies specified style edits one after another in the given order.

type SetStyleTextsCommand struct {
	params *SetStyleTextsParams
	result SetStyleTextsResult
	wg     sync.WaitGroup
	err    error
}

func NewSetStyleTextsCommand(params *SetStyleTextsParams) *SetStyleTextsCommand {
	return &SetStyleTextsCommand{
		params: params,
	}
}

func (cmd *SetStyleTextsCommand) Name() string {
	return "CSS.setStyleTexts"
}

func (cmd *SetStyleTextsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleTextsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetStyleTexts(params *SetStyleTextsParams, conn *hc.Conn) (result *SetStyleTextsResult, err error) {
	cmd := NewSetStyleTextsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetStyleTextsCB func(result *SetStyleTextsResult, err error)

// Applies specified style edits one after another in the given order.

type AsyncSetStyleTextsCommand struct {
	params *SetStyleTextsParams
	cb     SetStyleTextsCB
}

func NewAsyncSetStyleTextsCommand(params *SetStyleTextsParams, cb SetStyleTextsCB) *AsyncSetStyleTextsCommand {
	return &AsyncSetStyleTextsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetStyleTextsCommand) Name() string {
	return "CSS.setStyleTexts"
}

func (cmd *AsyncSetStyleTextsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleTextsCommand) Result() *SetStyleTextsResult {
	return &cmd.result
}

func (cmd *SetStyleTextsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetStyleTextsCommand) Done(data []byte, err error) {
	var result SetStyleTextsResult
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

type SetMediaTextParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
	Range        *SourceRange `json:"range"`
	Text         string       `json:"text"`
}

type SetMediaTextResult struct {
	Media *CSSMedia `json:"media"` // The resulting CSS media rule after modification.
}

// Modifies the rule selector.

type SetMediaTextCommand struct {
	params *SetMediaTextParams
	result SetMediaTextResult
	wg     sync.WaitGroup
	err    error
}

func NewSetMediaTextCommand(params *SetMediaTextParams) *SetMediaTextCommand {
	return &SetMediaTextCommand{
		params: params,
	}
}

func (cmd *SetMediaTextCommand) Name() string {
	return "CSS.setMediaText"
}

func (cmd *SetMediaTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMediaTextCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetMediaText(params *SetMediaTextParams, conn *hc.Conn) (result *SetMediaTextResult, err error) {
	cmd := NewSetMediaTextCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type SetMediaTextCB func(result *SetMediaTextResult, err error)

// Modifies the rule selector.

type AsyncSetMediaTextCommand struct {
	params *SetMediaTextParams
	cb     SetMediaTextCB
}

func NewAsyncSetMediaTextCommand(params *SetMediaTextParams, cb SetMediaTextCB) *AsyncSetMediaTextCommand {
	return &AsyncSetMediaTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetMediaTextCommand) Name() string {
	return "CSS.setMediaText"
}

func (cmd *AsyncSetMediaTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMediaTextCommand) Result() *SetMediaTextResult {
	return &cmd.result
}

func (cmd *SetMediaTextCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetMediaTextCommand) Done(data []byte, err error) {
	var result SetMediaTextResult
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

type CreateStyleSheetParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame where "via-inspector" stylesheet should be created.
}

type CreateStyleSheetResult struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // Identifier of the created "via-inspector" stylesheet.
}

// Creates a new special "via-inspector" stylesheet in the frame with given frameId.

type CreateStyleSheetCommand struct {
	params *CreateStyleSheetParams
	result CreateStyleSheetResult
	wg     sync.WaitGroup
	err    error
}

func NewCreateStyleSheetCommand(params *CreateStyleSheetParams) *CreateStyleSheetCommand {
	return &CreateStyleSheetCommand{
		params: params,
	}
}

func (cmd *CreateStyleSheetCommand) Name() string {
	return "CSS.createStyleSheet"
}

func (cmd *CreateStyleSheetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateStyleSheetCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func CreateStyleSheet(params *CreateStyleSheetParams, conn *hc.Conn) (result *CreateStyleSheetResult, err error) {
	cmd := NewCreateStyleSheetCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type CreateStyleSheetCB func(result *CreateStyleSheetResult, err error)

// Creates a new special "via-inspector" stylesheet in the frame with given frameId.

type AsyncCreateStyleSheetCommand struct {
	params *CreateStyleSheetParams
	cb     CreateStyleSheetCB
}

func NewAsyncCreateStyleSheetCommand(params *CreateStyleSheetParams, cb CreateStyleSheetCB) *AsyncCreateStyleSheetCommand {
	return &AsyncCreateStyleSheetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncCreateStyleSheetCommand) Name() string {
	return "CSS.createStyleSheet"
}

func (cmd *AsyncCreateStyleSheetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateStyleSheetCommand) Result() *CreateStyleSheetResult {
	return &cmd.result
}

func (cmd *CreateStyleSheetCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncCreateStyleSheetCommand) Done(data []byte, err error) {
	var result CreateStyleSheetResult
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

type AddRuleParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // The css style sheet identifier where a new rule should be inserted.
	RuleText     string       `json:"ruleText"`     // The text of a new rule.
	Location     *SourceRange `json:"location"`     // Text position of a new rule in the target style sheet.
}

type AddRuleResult struct {
	Rule *CSSRule `json:"rule"` // The newly created rule.
}

// Inserts a new rule with the given ruleText in a stylesheet with given styleSheetId, at the position specified by location.

type AddRuleCommand struct {
	params *AddRuleParams
	result AddRuleResult
	wg     sync.WaitGroup
	err    error
}

func NewAddRuleCommand(params *AddRuleParams) *AddRuleCommand {
	return &AddRuleCommand{
		params: params,
	}
}

func (cmd *AddRuleCommand) Name() string {
	return "CSS.addRule"
}

func (cmd *AddRuleCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddRuleCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func AddRule(params *AddRuleParams, conn *hc.Conn) (result *AddRuleResult, err error) {
	cmd := NewAddRuleCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type AddRuleCB func(result *AddRuleResult, err error)

// Inserts a new rule with the given ruleText in a stylesheet with given styleSheetId, at the position specified by location.

type AsyncAddRuleCommand struct {
	params *AddRuleParams
	cb     AddRuleCB
}

func NewAsyncAddRuleCommand(params *AddRuleParams, cb AddRuleCB) *AsyncAddRuleCommand {
	return &AsyncAddRuleCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncAddRuleCommand) Name() string {
	return "CSS.addRule"
}

func (cmd *AsyncAddRuleCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddRuleCommand) Result() *AddRuleResult {
	return &cmd.result
}

func (cmd *AddRuleCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncAddRuleCommand) Done(data []byte, err error) {
	var result AddRuleResult
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

type ForcePseudoStateParams struct {
	NodeId              *NodeId  `json:"nodeId"`              // The element id for which to force the pseudo state.
	ForcedPseudoClasses []string `json:"forcedPseudoClasses"` // Element pseudo classes to force when computing the element's style.
}

// Ensures that the given node will have specified pseudo-classes whenever its style is computed by the browser.

type ForcePseudoStateCommand struct {
	params *ForcePseudoStateParams
	wg     sync.WaitGroup
	err    error
}

func NewForcePseudoStateCommand(params *ForcePseudoStateParams) *ForcePseudoStateCommand {
	return &ForcePseudoStateCommand{
		params: params,
	}
}

func (cmd *ForcePseudoStateCommand) Name() string {
	return "CSS.forcePseudoState"
}

func (cmd *ForcePseudoStateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForcePseudoStateCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func ForcePseudoState(params *ForcePseudoStateParams, conn *hc.Conn) (err error) {
	cmd := NewForcePseudoStateCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type ForcePseudoStateCB func(err error)

// Ensures that the given node will have specified pseudo-classes whenever its style is computed by the browser.

type AsyncForcePseudoStateCommand struct {
	params *ForcePseudoStateParams
	cb     ForcePseudoStateCB
}

func NewAsyncForcePseudoStateCommand(params *ForcePseudoStateParams, cb ForcePseudoStateCB) *AsyncForcePseudoStateCommand {
	return &AsyncForcePseudoStateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncForcePseudoStateCommand) Name() string {
	return "CSS.forcePseudoState"
}

func (cmd *AsyncForcePseudoStateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForcePseudoStateCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncForcePseudoStateCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetMediaQueriesResult struct {
	Medias []*CSSMedia `json:"medias"`
}

// Returns all media queries parsed by the rendering engine.
// @experimental
type GetMediaQueriesCommand struct {
	result GetMediaQueriesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetMediaQueriesCommand() *GetMediaQueriesCommand {
	return &GetMediaQueriesCommand{}
}

func (cmd *GetMediaQueriesCommand) Name() string {
	return "CSS.getMediaQueries"
}

func (cmd *GetMediaQueriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetMediaQueriesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetMediaQueries(conn *hc.Conn) (result *GetMediaQueriesResult, err error) {
	cmd := NewGetMediaQueriesCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetMediaQueriesCB func(result *GetMediaQueriesResult, err error)

// Returns all media queries parsed by the rendering engine.
// @experimental
type AsyncGetMediaQueriesCommand struct {
	cb GetMediaQueriesCB
}

func NewAsyncGetMediaQueriesCommand(cb GetMediaQueriesCB) *AsyncGetMediaQueriesCommand {
	return &AsyncGetMediaQueriesCommand{
		cb: cb,
	}
}

func (cmd *AsyncGetMediaQueriesCommand) Name() string {
	return "CSS.getMediaQueries"
}

func (cmd *AsyncGetMediaQueriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetMediaQueriesCommand) Result() *GetMediaQueriesResult {
	return &cmd.result
}

func (cmd *GetMediaQueriesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetMediaQueriesCommand) Done(data []byte, err error) {
	var result GetMediaQueriesResult
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

type SetEffectivePropertyValueForNodeParams struct {
	NodeId       *NodeId `json:"nodeId"` // The element id for which to set property.
	PropertyName string  `json:"propertyName"`
	Value        string  `json:"value"`
}

// Find a rule with the given active property for the given node and set the new value for this property
// @experimental
type SetEffectivePropertyValueForNodeCommand struct {
	params *SetEffectivePropertyValueForNodeParams
	wg     sync.WaitGroup
	err    error
}

func NewSetEffectivePropertyValueForNodeCommand(params *SetEffectivePropertyValueForNodeParams) *SetEffectivePropertyValueForNodeCommand {
	return &SetEffectivePropertyValueForNodeCommand{
		params: params,
	}
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Name() string {
	return "CSS.setEffectivePropertyValueForNode"
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func SetEffectivePropertyValueForNode(params *SetEffectivePropertyValueForNodeParams, conn *hc.Conn) (err error) {
	cmd := NewSetEffectivePropertyValueForNodeCommand(params)
	cmd.Run(conn)
	return cmd.err
}

type SetEffectivePropertyValueForNodeCB func(err error)

// Find a rule with the given active property for the given node and set the new value for this property
// @experimental
type AsyncSetEffectivePropertyValueForNodeCommand struct {
	params *SetEffectivePropertyValueForNodeParams
	cb     SetEffectivePropertyValueForNodeCB
}

func NewAsyncSetEffectivePropertyValueForNodeCommand(params *SetEffectivePropertyValueForNodeParams, cb SetEffectivePropertyValueForNodeCB) *AsyncSetEffectivePropertyValueForNodeCommand {
	return &AsyncSetEffectivePropertyValueForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncSetEffectivePropertyValueForNodeCommand) Name() string {
	return "CSS.setEffectivePropertyValueForNode"
}

func (cmd *AsyncSetEffectivePropertyValueForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncSetEffectivePropertyValueForNodeCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type GetBackgroundColorsParams struct {
	NodeId *NodeId `json:"nodeId"` // Id of the node to get background colors for.
}

type GetBackgroundColorsResult struct {
	BackgroundColors []string `json:"backgroundColors"` // The range of background colors behind this element, if it contains any visible text. If no visible text is present, this will be undefined. In the case of a flat background color, this will consist of simply that color. In the case of a gradient, this will consist of each of the color stops. For anything more complicated, this will be an empty array. Images will be ignored (as if the image had failed to load).
}

// @experimental
type GetBackgroundColorsCommand struct {
	params *GetBackgroundColorsParams
	result GetBackgroundColorsResult
	wg     sync.WaitGroup
	err    error
}

func NewGetBackgroundColorsCommand(params *GetBackgroundColorsParams) *GetBackgroundColorsCommand {
	return &GetBackgroundColorsCommand{
		params: params,
	}
}

func (cmd *GetBackgroundColorsCommand) Name() string {
	return "CSS.getBackgroundColors"
}

func (cmd *GetBackgroundColorsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBackgroundColorsCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetBackgroundColors(params *GetBackgroundColorsParams, conn *hc.Conn) (result *GetBackgroundColorsResult, err error) {
	cmd := NewGetBackgroundColorsCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetBackgroundColorsCB func(result *GetBackgroundColorsResult, err error)

// @experimental
type AsyncGetBackgroundColorsCommand struct {
	params *GetBackgroundColorsParams
	cb     GetBackgroundColorsCB
}

func NewAsyncGetBackgroundColorsCommand(params *GetBackgroundColorsParams, cb GetBackgroundColorsCB) *AsyncGetBackgroundColorsCommand {
	return &AsyncGetBackgroundColorsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetBackgroundColorsCommand) Name() string {
	return "CSS.getBackgroundColors"
}

func (cmd *AsyncGetBackgroundColorsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBackgroundColorsCommand) Result() *GetBackgroundColorsResult {
	return &cmd.result
}

func (cmd *GetBackgroundColorsCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetBackgroundColorsCommand) Done(data []byte, err error) {
	var result GetBackgroundColorsResult
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

type GetLayoutTreeAndStylesParams struct {
	ComputedStyleWhitelist []string `json:"computedStyleWhitelist"` // Whitelist of computed styles to return.
}

type GetLayoutTreeAndStylesResult struct {
	LayoutTreeNodes []*LayoutTreeNode `json:"layoutTreeNodes"`
	ComputedStyles  []*ComputedStyle  `json:"computedStyles"`
}

// For the main document and any content documents, return the LayoutTreeNodes and a whitelisted subset of the computed style. It only returns pushed nodes, on way to pull all nodes is to call DOM.getDocument with a depth of -1.
// @experimental
type GetLayoutTreeAndStylesCommand struct {
	params *GetLayoutTreeAndStylesParams
	result GetLayoutTreeAndStylesResult
	wg     sync.WaitGroup
	err    error
}

func NewGetLayoutTreeAndStylesCommand(params *GetLayoutTreeAndStylesParams) *GetLayoutTreeAndStylesCommand {
	return &GetLayoutTreeAndStylesCommand{
		params: params,
	}
}

func (cmd *GetLayoutTreeAndStylesCommand) Name() string {
	return "CSS.getLayoutTreeAndStyles"
}

func (cmd *GetLayoutTreeAndStylesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetLayoutTreeAndStylesCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func GetLayoutTreeAndStyles(params *GetLayoutTreeAndStylesParams, conn *hc.Conn) (result *GetLayoutTreeAndStylesResult, err error) {
	cmd := NewGetLayoutTreeAndStylesCommand(params)
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type GetLayoutTreeAndStylesCB func(result *GetLayoutTreeAndStylesResult, err error)

// For the main document and any content documents, return the LayoutTreeNodes and a whitelisted subset of the computed style. It only returns pushed nodes, on way to pull all nodes is to call DOM.getDocument with a depth of -1.
// @experimental
type AsyncGetLayoutTreeAndStylesCommand struct {
	params *GetLayoutTreeAndStylesParams
	cb     GetLayoutTreeAndStylesCB
}

func NewAsyncGetLayoutTreeAndStylesCommand(params *GetLayoutTreeAndStylesParams, cb GetLayoutTreeAndStylesCB) *AsyncGetLayoutTreeAndStylesCommand {
	return &AsyncGetLayoutTreeAndStylesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AsyncGetLayoutTreeAndStylesCommand) Name() string {
	return "CSS.getLayoutTreeAndStyles"
}

func (cmd *AsyncGetLayoutTreeAndStylesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetLayoutTreeAndStylesCommand) Result() *GetLayoutTreeAndStylesResult {
	return &cmd.result
}

func (cmd *GetLayoutTreeAndStylesCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncGetLayoutTreeAndStylesCommand) Done(data []byte, err error) {
	var result GetLayoutTreeAndStylesResult
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

// Enables the selector recording.
// @experimental
type StartRuleUsageTrackingCommand struct {
	wg  sync.WaitGroup
	err error
}

func NewStartRuleUsageTrackingCommand() *StartRuleUsageTrackingCommand {
	return &StartRuleUsageTrackingCommand{}
}

func (cmd *StartRuleUsageTrackingCommand) Name() string {
	return "CSS.startRuleUsageTracking"
}

func (cmd *StartRuleUsageTrackingCommand) Params() interface{} {
	return nil
}

func (cmd *StartRuleUsageTrackingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StartRuleUsageTracking(conn *hc.Conn) (err error) {
	cmd := NewStartRuleUsageTrackingCommand()
	cmd.Run(conn)
	return cmd.err
}

type StartRuleUsageTrackingCB func(err error)

// Enables the selector recording.
// @experimental
type AsyncStartRuleUsageTrackingCommand struct {
	cb StartRuleUsageTrackingCB
}

func NewAsyncStartRuleUsageTrackingCommand(cb StartRuleUsageTrackingCB) *AsyncStartRuleUsageTrackingCommand {
	return &AsyncStartRuleUsageTrackingCommand{
		cb: cb,
	}
}

func (cmd *AsyncStartRuleUsageTrackingCommand) Name() string {
	return "CSS.startRuleUsageTracking"
}

func (cmd *AsyncStartRuleUsageTrackingCommand) Params() interface{} {
	return nil
}

func (cmd *StartRuleUsageTrackingCommand) Done(data []byte, err error) {
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStartRuleUsageTrackingCommand) Done(data []byte, err error) {
	cmd.cb(err)
}

type StopRuleUsageTrackingResult struct {
	RuleUsage []*RuleUsage `json:"ruleUsage"`
}

// The list of rules with an indication of whether these were used
// @experimental
type StopRuleUsageTrackingCommand struct {
	result StopRuleUsageTrackingResult
	wg     sync.WaitGroup
	err    error
}

func NewStopRuleUsageTrackingCommand() *StopRuleUsageTrackingCommand {
	return &StopRuleUsageTrackingCommand{}
}

func (cmd *StopRuleUsageTrackingCommand) Name() string {
	return "CSS.stopRuleUsageTracking"
}

func (cmd *StopRuleUsageTrackingCommand) Params() interface{} {
	return nil
}

func (cmd *StopRuleUsageTrackingCommand) Run(conn *hc.Conn) error {
	cmd.wg.Add(1)
	conn.SendCommand(cmd)
	cmd.wg.Wait()
	return cmd.err
}

func StopRuleUsageTracking(conn *hc.Conn) (result *StopRuleUsageTrackingResult, err error) {
	cmd := NewStopRuleUsageTrackingCommand()
	cmd.Run(conn)
	return &cmd.result, cmd.err
}

type StopRuleUsageTrackingCB func(result *StopRuleUsageTrackingResult, err error)

// The list of rules with an indication of whether these were used
// @experimental
type AsyncStopRuleUsageTrackingCommand struct {
	cb StopRuleUsageTrackingCB
}

func NewAsyncStopRuleUsageTrackingCommand(cb StopRuleUsageTrackingCB) *AsyncStopRuleUsageTrackingCommand {
	return &AsyncStopRuleUsageTrackingCommand{
		cb: cb,
	}
}

func (cmd *AsyncStopRuleUsageTrackingCommand) Name() string {
	return "CSS.stopRuleUsageTracking"
}

func (cmd *AsyncStopRuleUsageTrackingCommand) Params() interface{} {
	return nil
}

func (cmd *StopRuleUsageTrackingCommand) Result() *StopRuleUsageTrackingResult {
	return &cmd.result
}

func (cmd *StopRuleUsageTrackingCommand) Done(data []byte, err error) {
	if err == nil {
		err = json.Unmarshal(data, &cmd.result)
	}
	cmd.err = err
	cmd.wg.Done()
}

func (cmd *AsyncStopRuleUsageTrackingCommand) Done(data []byte, err error) {
	var result StopRuleUsageTrackingResult
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

// Fires whenever a MediaQuery result changes (for example, after a browser window has been resized.) The current implementation considers only viewport-dependent media features.

type MediaQueryResultChangedEvent struct {
}

func OnMediaQueryResultChanged(conn *hc.Conn, cb func(evt *MediaQueryResultChangedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &MediaQueryResultChangedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("CSS.mediaQueryResultChanged", sink)
}

// Fires whenever a web font gets loaded.

type FontsUpdatedEvent struct {
}

func OnFontsUpdated(conn *hc.Conn, cb func(evt *FontsUpdatedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &FontsUpdatedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("CSS.fontsUpdated", sink)
}

// Fired whenever a stylesheet is changed as a result of the client operation.

type StyleSheetChangedEvent struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

func OnStyleSheetChanged(conn *hc.Conn, cb func(evt *StyleSheetChangedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &StyleSheetChangedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("CSS.styleSheetChanged", sink)
}

// Fired whenever an active document stylesheet is added.

type StyleSheetAddedEvent struct {
	Header *CSSStyleSheetHeader `json:"header"` // Added stylesheet metainfo.
}

func OnStyleSheetAdded(conn *hc.Conn, cb func(evt *StyleSheetAddedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &StyleSheetAddedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("CSS.styleSheetAdded", sink)
}

// Fired whenever an active document stylesheet is removed.

type StyleSheetRemovedEvent struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // Identifier of the removed stylesheet.
}

func OnStyleSheetRemoved(conn *hc.Conn, cb func(evt *StyleSheetRemovedEvent)) {
	sink := hc.FuncToEventSink(func(name string, params []byte) {
		evt := &StyleSheetRemovedEvent{}
		if err := json.Unmarshal(params, evt); err != nil {
			logging.Vlog(-1, err)
		} else {
			cb(evt)
		}
	})
	conn.AddEventSink("CSS.styleSheetRemoved", sink)
}
