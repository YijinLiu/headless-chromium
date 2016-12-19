package protocol

import (
	"encoding/json"
	"github.com/yijinliu/algo-lib/go/src/logging"
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
	InlineStyle     *CSSStyle    `json:"inlineStyle"`     // The ancestor node's inline style, if any, in the style inheritance chain.
	MatchedCSSRules []*RuleMatch `json:"matchedCSSRules"` // Matches of CSS rules matching the ancestor node in the style inheritance chain.
}

// Match data for a CSS rule.
type RuleMatch struct {
	Rule              *CSSRule `json:"rule"`              // CSS rule in the match.
	MatchingSelectors []int    `json:"matchingSelectors"` // Matching selector indices in the rule's selectorList selectors (0-based).
}

// Data for a simple selector (these are delimited by commas in a selector list).
type Value struct {
	Text  string       `json:"text"`  // Value text.
	Range *SourceRange `json:"range"` // Value range in the underlying resource (if available).
}

// Selector list data.
type SelectorList struct {
	Selectors []*Value `json:"selectors"` // Selectors in the list.
	Text      string   `json:"text"`      // Rule selector text.
}

// CSS stylesheet metainformation.
type CSSStyleSheetHeader struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId"` // The stylesheet identifier.
	FrameId      *FrameId         `json:"frameId"`      // Owner frame identifier.
	SourceURL    string           `json:"sourceURL"`    // Stylesheet resource URL.
	SourceMapURL string           `json:"sourceMapURL"` // URL of source map associated with the stylesheet (if any).
	Origin       StyleSheetOrigin `json:"origin"`       // Stylesheet origin.
	Title        string           `json:"title"`        // Stylesheet title.
	OwnerNode    *BackendNodeId   `json:"ownerNode"`    // The backend id for the owner node of the stylesheet.
	Disabled     bool             `json:"disabled"`     // Denotes whether the stylesheet is disabled.
	HasSourceURL bool             `json:"hasSourceURL"` // Whether the sourceURL field value comes from the sourceURL comment.
	IsInline     bool             `json:"isInline"`     // Whether this stylesheet is created for STYLE tag by parser. This flag is not set for document.written STYLE tags.
	StartLine    int              `json:"startLine"`    // Line offset of the stylesheet within the resource (zero based).
	StartColumn  int              `json:"startColumn"`  // Column offset of the stylesheet within the resource (zero based).
}

// CSS rule representation.
type CSSRule struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	SelectorList *SelectorList    `json:"selectorList"` // Rule selector data.
	Origin       StyleSheetOrigin `json:"origin"`       // Parent stylesheet's origin.
	Style        *CSSStyle        `json:"style"`        // Associated style declaration.
	Media        []*CSSMedia      `json:"media"`        // Media list array (for rules involving media queries). The array enumerates media queries starting with the innermost one, going outwards.
}

// Text range within a resource. All numbers are zero-based.
type SourceRange struct {
	StartLine   int `json:"startLine"`   // Start line of range.
	StartColumn int `json:"startColumn"` // Start column of range (inclusive).
	EndLine     int `json:"endLine"`     // End line of range
	EndColumn   int `json:"endColumn"`   // End column of range (exclusive).
}

type ShorthandEntry struct {
	Name      string `json:"name"`      // Shorthand name.
	Value     string `json:"value"`     // Shorthand value.
	Important bool   `json:"important"` // Whether the property has "!important" annotation (implies false if absent).
}

type CSSComputedStyleProperty struct {
	Name  string `json:"name"`  // Computed style property name.
	Value string `json:"value"` // Computed style property value.
}

// CSS style representation.
type CSSStyle struct {
	StyleSheetId     StyleSheetId      `json:"styleSheetId"`     // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	CssProperties    []*CSSProperty    `json:"cssProperties"`    // CSS properties in the style.
	ShorthandEntries []*ShorthandEntry `json:"shorthandEntries"` // Computed values for all shorthands found in the style.
	CssText          string            `json:"cssText"`          // Style declaration text (if available).
	Range            *SourceRange      `json:"range"`            // Style declaration range in the enclosing stylesheet (if available).
}

// CSS property declaration data.
type CSSProperty struct {
	Name      string       `json:"name"`      // The property name.
	Value     string       `json:"value"`     // The property value.
	Important bool         `json:"important"` // Whether the property has "!important" annotation (implies false if absent).
	Implicit  bool         `json:"implicit"`  // Whether the property is implicit (implies false if absent).
	Text      string       `json:"text"`      // The full property text as specified in the style.
	ParsedOk  bool         `json:"parsedOk"`  // Whether the property is understood by the browser (implies true if absent).
	Disabled  bool         `json:"disabled"`  // Whether the property is disabled by the user (present for source-based properties only).
	Range     *SourceRange `json:"range"`     // The entire property range in the enclosing style declaration (if available).
}

// CSS media rule descriptor.
type CSSMedia struct {
	Text         string        `json:"text"`         // Media query text.
	Source       string        `json:"source"`       // Source of the media query: "mediaRule" if specified by a @media rule, "importRule" if specified by an @import rule, "linkedSheet" if specified by a "media" attribute in a linked stylesheet's LINK tag, "inlineSheet" if specified by a "media" attribute in an inline stylesheet's STYLE tag.
	SourceURL    string        `json:"sourceURL"`    // URL of the document containing the media query description.
	Range        *SourceRange  `json:"range"`        // The associated rule (@media or @import) header range in the enclosing stylesheet (if available).
	StyleSheetId StyleSheetId  `json:"styleSheetId"` // Identifier of the stylesheet containing this object (if exists).
	MediaList    []*MediaQuery `json:"mediaList"`    // Array of media queries.
}

// Media query descriptor.
type MediaQuery struct {
	Expressions []*MediaQueryExpression `json:"expressions"` // Array of media query expressions.
	Active      bool                    `json:"active"`      // Whether the media query condition is satisfied.
}

// Media query expression descriptor.
type MediaQueryExpression struct {
	Value          int          `json:"value"`          // Media query expression value.
	Unit           string       `json:"unit"`           // Media query expression units.
	Feature        string       `json:"feature"`        // Media query expression feature.
	ValueRange     *SourceRange `json:"valueRange"`     // The associated range of the value text in the enclosing stylesheet (if available).
	ComputedLength int          `json:"computedLength"` // Computed length of media query expression (if applicable).
}

// Information about amount of glyphs that were rendered with given font.
type PlatformFontUsage struct {
	FamilyName   string `json:"familyName"`   // Font's family name reported by platform.
	IsCustomFont bool   `json:"isCustomFont"` // Indicates if the font was downloaded or resolved locally.
	GlyphCount   int    `json:"glyphCount"`   // Amount of glyphs that were rendered with this font.
}

// CSS keyframes rule representation.
type CSSKeyframesRule struct {
	AnimationName *Value             `json:"animationName"` // Animation name.
	Keyframes     []*CSSKeyframeRule `json:"keyframes"`     // List of keyframes.
}

// CSS keyframe rule representation.
type CSSKeyframeRule struct {
	StyleSheetId StyleSheetId     `json:"styleSheetId"` // The css style sheet identifier (absent for user agent stylesheet and user-specified stylesheet rules) this rule came from.
	Origin       StyleSheetOrigin `json:"origin"`       // Parent stylesheet's origin.
	KeyText      *Value           `json:"keyText"`      // Associated key text.
	Style        *CSSStyle        `json:"style"`        // Associated style declaration.
}

// A descriptor of operation to mutate style declaration text.
type StyleDeclarationEdit struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // The css style sheet identifier.
	Range        *SourceRange `json:"range"`        // The range of the style text in the enclosing stylesheet.
	Text         string       `json:"text"`         // New style text.
}

type CSSEnableCB func(err error)

// Enables the CSS agent for the given page. Clients should not assume that the CSS agent has been enabled until the result of this command is received.
type CSSEnableCommand struct {
	cb CSSEnableCB
}

func NewCSSEnableCommand(cb CSSEnableCB) *CSSEnableCommand {
	return &CSSEnableCommand{
		cb: cb,
	}
}

func (cmd *CSSEnableCommand) Name() string {
	return "CSS.enable"
}

func (cmd *CSSEnableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSEnableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type CSSDisableCB func(err error)

// Disables the CSS agent for the given page.
type CSSDisableCommand struct {
	cb CSSDisableCB
}

func NewCSSDisableCommand(cb CSSDisableCB) *CSSDisableCommand {
	return &CSSDisableCommand{
		cb: cb,
	}
}

func (cmd *CSSDisableCommand) Name() string {
	return "CSS.disable"
}

func (cmd *CSSDisableCommand) Params() interface{} {
	return nil
}

func (cmd *CSSDisableCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
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

type GetMatchedStylesForNodeCB func(result *GetMatchedStylesForNodeResult, err error)

// Returns requested styles for a DOM node identified by nodeId.
type GetMatchedStylesForNodeCommand struct {
	params *GetMatchedStylesForNodeParams
	cb     GetMatchedStylesForNodeCB
}

func NewGetMatchedStylesForNodeCommand(params *GetMatchedStylesForNodeParams, cb GetMatchedStylesForNodeCB) *GetMatchedStylesForNodeCommand {
	return &GetMatchedStylesForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetMatchedStylesForNodeCommand) Name() string {
	return "CSS.getMatchedStylesForNode"
}

func (cmd *GetMatchedStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetMatchedStylesForNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetMatchedStylesForNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetInlineStylesForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetInlineStylesForNodeResult struct {
	InlineStyle     *CSSStyle `json:"inlineStyle"`     // Inline style for the specified DOM node.
	AttributesStyle *CSSStyle `json:"attributesStyle"` // Attribute-defined element style (e.g. resulting from "width=20 height=100%").
}

type GetInlineStylesForNodeCB func(result *GetInlineStylesForNodeResult, err error)

// Returns the styles defined inline (explicitly in the "style" attribute and implicitly, using DOM attributes) for a DOM node identified by nodeId.
type GetInlineStylesForNodeCommand struct {
	params *GetInlineStylesForNodeParams
	cb     GetInlineStylesForNodeCB
}

func NewGetInlineStylesForNodeCommand(params *GetInlineStylesForNodeParams, cb GetInlineStylesForNodeCB) *GetInlineStylesForNodeCommand {
	return &GetInlineStylesForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetInlineStylesForNodeCommand) Name() string {
	return "CSS.getInlineStylesForNode"
}

func (cmd *GetInlineStylesForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetInlineStylesForNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetInlineStylesForNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetComputedStyleForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetComputedStyleForNodeResult struct {
	ComputedStyle []*CSSComputedStyleProperty `json:"computedStyle"` // Computed style for the specified DOM node.
}

type GetComputedStyleForNodeCB func(result *GetComputedStyleForNodeResult, err error)

// Returns the computed style for a DOM node identified by nodeId.
type GetComputedStyleForNodeCommand struct {
	params *GetComputedStyleForNodeParams
	cb     GetComputedStyleForNodeCB
}

func NewGetComputedStyleForNodeCommand(params *GetComputedStyleForNodeParams, cb GetComputedStyleForNodeCB) *GetComputedStyleForNodeCommand {
	return &GetComputedStyleForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetComputedStyleForNodeCommand) Name() string {
	return "CSS.getComputedStyleForNode"
}

func (cmd *GetComputedStyleForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetComputedStyleForNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetComputedStyleForNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetPlatformFontsForNodeParams struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetPlatformFontsForNodeResult struct {
	Fonts []*PlatformFontUsage `json:"fonts"` // Usage statistics for every employed platform font.
}

type GetPlatformFontsForNodeCB func(result *GetPlatformFontsForNodeResult, err error)

// Requests information about platform fonts which we used to render child TextNodes in the given node.
type GetPlatformFontsForNodeCommand struct {
	params *GetPlatformFontsForNodeParams
	cb     GetPlatformFontsForNodeCB
}

func NewGetPlatformFontsForNodeCommand(params *GetPlatformFontsForNodeParams, cb GetPlatformFontsForNodeCB) *GetPlatformFontsForNodeCommand {
	return &GetPlatformFontsForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetPlatformFontsForNodeCommand) Name() string {
	return "CSS.getPlatformFontsForNode"
}

func (cmd *GetPlatformFontsForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetPlatformFontsForNodeCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetPlatformFontsForNodeResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type GetStyleSheetTextParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

type GetStyleSheetTextResult struct {
	Text string `json:"text"` // The stylesheet text.
}

type GetStyleSheetTextCB func(result *GetStyleSheetTextResult, err error)

// Returns the current textual content and the URL for a stylesheet.
type GetStyleSheetTextCommand struct {
	params *GetStyleSheetTextParams
	cb     GetStyleSheetTextCB
}

func NewGetStyleSheetTextCommand(params *GetStyleSheetTextParams, cb GetStyleSheetTextCB) *GetStyleSheetTextCommand {
	return &GetStyleSheetTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetStyleSheetTextCommand) Name() string {
	return "CSS.getStyleSheetText"
}

func (cmd *GetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetStyleSheetTextCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetStyleSheetTextResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CollectClassNamesParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

type CollectClassNamesResult struct {
	ClassNames []string `json:"classNames"` // Class name list.
}

type CollectClassNamesCB func(result *CollectClassNamesResult, err error)

// Returns all class names from specified stylesheet.
type CollectClassNamesCommand struct {
	params *CollectClassNamesParams
	cb     CollectClassNamesCB
}

func NewCollectClassNamesCommand(params *CollectClassNamesParams, cb CollectClassNamesCB) *CollectClassNamesCommand {
	return &CollectClassNamesCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CollectClassNamesCommand) Name() string {
	return "CSS.collectClassNames"
}

func (cmd *CollectClassNamesCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CollectClassNamesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CollectClassNamesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetStyleSheetTextParams struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
	Text         string       `json:"text"`
}

type SetStyleSheetTextResult struct {
	SourceMapURL string `json:"sourceMapURL"` // URL of source map associated with script (if any).
}

type SetStyleSheetTextCB func(result *SetStyleSheetTextResult, err error)

// Sets the new stylesheet text.
type SetStyleSheetTextCommand struct {
	params *SetStyleSheetTextParams
	cb     SetStyleSheetTextCB
}

func NewSetStyleSheetTextCommand(params *SetStyleSheetTextParams, cb SetStyleSheetTextCB) *SetStyleSheetTextCommand {
	return &SetStyleSheetTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetStyleSheetTextCommand) Name() string {
	return "CSS.setStyleSheetText"
}

func (cmd *SetStyleSheetTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleSheetTextCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetStyleSheetTextResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type SetRuleSelectorCB func(result *SetRuleSelectorResult, err error)

// Modifies the rule selector.
type SetRuleSelectorCommand struct {
	params *SetRuleSelectorParams
	cb     SetRuleSelectorCB
}

func NewSetRuleSelectorCommand(params *SetRuleSelectorParams, cb SetRuleSelectorCB) *SetRuleSelectorCommand {
	return &SetRuleSelectorCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetRuleSelectorCommand) Name() string {
	return "CSS.setRuleSelector"
}

func (cmd *SetRuleSelectorCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetRuleSelectorCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetRuleSelectorResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type SetKeyframeKeyCB func(result *SetKeyframeKeyResult, err error)

// Modifies the keyframe rule key text.
type SetKeyframeKeyCommand struct {
	params *SetKeyframeKeyParams
	cb     SetKeyframeKeyCB
}

func NewSetKeyframeKeyCommand(params *SetKeyframeKeyParams, cb SetKeyframeKeyCB) *SetKeyframeKeyCommand {
	return &SetKeyframeKeyCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetKeyframeKeyCommand) Name() string {
	return "CSS.setKeyframeKey"
}

func (cmd *SetKeyframeKeyCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetKeyframeKeyCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetKeyframeKeyResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetStyleTextsParams struct {
	Edits []*StyleDeclarationEdit `json:"edits"`
}

type SetStyleTextsResult struct {
	Styles []*CSSStyle `json:"styles"` // The resulting styles after modification.
}

type SetStyleTextsCB func(result *SetStyleTextsResult, err error)

// Applies specified style edits one after another in the given order.
type SetStyleTextsCommand struct {
	params *SetStyleTextsParams
	cb     SetStyleTextsCB
}

func NewSetStyleTextsCommand(params *SetStyleTextsParams, cb SetStyleTextsCB) *SetStyleTextsCommand {
	return &SetStyleTextsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetStyleTextsCommand) Name() string {
	return "CSS.setStyleTexts"
}

func (cmd *SetStyleTextsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetStyleTextsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetStyleTextsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type SetMediaTextCB func(result *SetMediaTextResult, err error)

// Modifies the rule selector.
type SetMediaTextCommand struct {
	params *SetMediaTextParams
	cb     SetMediaTextCB
}

func NewSetMediaTextCommand(params *SetMediaTextParams, cb SetMediaTextCB) *SetMediaTextCommand {
	return &SetMediaTextCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetMediaTextCommand) Name() string {
	return "CSS.setMediaText"
}

func (cmd *SetMediaTextCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetMediaTextCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj SetMediaTextResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type CreateStyleSheetParams struct {
	FrameId *FrameId `json:"frameId"` // Identifier of the frame where "via-inspector" stylesheet should be created.
}

type CreateStyleSheetResult struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // Identifier of the created "via-inspector" stylesheet.
}

type CreateStyleSheetCB func(result *CreateStyleSheetResult, err error)

// Creates a new special "via-inspector" stylesheet in the frame with given frameId.
type CreateStyleSheetCommand struct {
	params *CreateStyleSheetParams
	cb     CreateStyleSheetCB
}

func NewCreateStyleSheetCommand(params *CreateStyleSheetParams, cb CreateStyleSheetCB) *CreateStyleSheetCommand {
	return &CreateStyleSheetCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *CreateStyleSheetCommand) Name() string {
	return "CSS.createStyleSheet"
}

func (cmd *CreateStyleSheetCommand) Params() interface{} {
	return cmd.params
}

func (cmd *CreateStyleSheetCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj CreateStyleSheetResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
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

type AddRuleCB func(result *AddRuleResult, err error)

// Inserts a new rule with the given ruleText in a stylesheet with given styleSheetId, at the position specified by location.
type AddRuleCommand struct {
	params *AddRuleParams
	cb     AddRuleCB
}

func NewAddRuleCommand(params *AddRuleParams, cb AddRuleCB) *AddRuleCommand {
	return &AddRuleCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *AddRuleCommand) Name() string {
	return "CSS.addRule"
}

func (cmd *AddRuleCommand) Params() interface{} {
	return cmd.params
}

func (cmd *AddRuleCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj AddRuleResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type ForcePseudoStateParams struct {
	NodeId              *NodeId  `json:"nodeId"`              // The element id for which to force the pseudo state.
	ForcedPseudoClasses []string `json:"forcedPseudoClasses"` // Element pseudo classes to force when computing the element's style.
}

type ForcePseudoStateCB func(err error)

// Ensures that the given node will have specified pseudo-classes whenever its style is computed by the browser.
type ForcePseudoStateCommand struct {
	params *ForcePseudoStateParams
	cb     ForcePseudoStateCB
}

func NewForcePseudoStateCommand(params *ForcePseudoStateParams, cb ForcePseudoStateCB) *ForcePseudoStateCommand {
	return &ForcePseudoStateCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *ForcePseudoStateCommand) Name() string {
	return "CSS.forcePseudoState"
}

func (cmd *ForcePseudoStateCommand) Params() interface{} {
	return cmd.params
}

func (cmd *ForcePseudoStateCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetMediaQueriesResult struct {
	Medias []*CSSMedia `json:"medias"`
}

type GetMediaQueriesCB func(result *GetMediaQueriesResult, err error)

// Returns all media queries parsed by the rendering engine.
type GetMediaQueriesCommand struct {
	cb GetMediaQueriesCB
}

func NewGetMediaQueriesCommand(cb GetMediaQueriesCB) *GetMediaQueriesCommand {
	return &GetMediaQueriesCommand{
		cb: cb,
	}
}

func (cmd *GetMediaQueriesCommand) Name() string {
	return "CSS.getMediaQueries"
}

func (cmd *GetMediaQueriesCommand) Params() interface{} {
	return nil
}

func (cmd *GetMediaQueriesCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetMediaQueriesResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type SetEffectivePropertyValueForNodeParams struct {
	NodeId       *NodeId `json:"nodeId"` // The element id for which to set property.
	PropertyName string  `json:"propertyName"`
	Value        string  `json:"value"`
}

type SetEffectivePropertyValueForNodeCB func(err error)

// Find a rule with the given active property for the given node and set the new value for this property
type SetEffectivePropertyValueForNodeCommand struct {
	params *SetEffectivePropertyValueForNodeParams
	cb     SetEffectivePropertyValueForNodeCB
}

func NewSetEffectivePropertyValueForNodeCommand(params *SetEffectivePropertyValueForNodeParams, cb SetEffectivePropertyValueForNodeCB) *SetEffectivePropertyValueForNodeCommand {
	return &SetEffectivePropertyValueForNodeCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Name() string {
	return "CSS.setEffectivePropertyValueForNode"
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Params() interface{} {
	return cmd.params
}

func (cmd *SetEffectivePropertyValueForNodeCommand) Done(result []byte, err error) {
	if cmd.cb != nil {
		cmd.cb(err)
	}
}

type GetBackgroundColorsParams struct {
	NodeId *NodeId `json:"nodeId"` // Id of the node to get background colors for.
}

type GetBackgroundColorsResult struct {
	BackgroundColors []string `json:"backgroundColors"` // The range of background colors behind this element, if it contains any visible text. If no visible text is present, this will be undefined. In the case of a flat background color, this will consist of simply that color. In the case of a gradient, this will consist of each of the color stops. For anything more complicated, this will be an empty array. Images will be ignored (as if the image had failed to load).
}

type GetBackgroundColorsCB func(result *GetBackgroundColorsResult, err error)

type GetBackgroundColorsCommand struct {
	params *GetBackgroundColorsParams
	cb     GetBackgroundColorsCB
}

func NewGetBackgroundColorsCommand(params *GetBackgroundColorsParams, cb GetBackgroundColorsCB) *GetBackgroundColorsCommand {
	return &GetBackgroundColorsCommand{
		params: params,
		cb:     cb,
	}
}

func (cmd *GetBackgroundColorsCommand) Name() string {
	return "CSS.getBackgroundColors"
}

func (cmd *GetBackgroundColorsCommand) Params() interface{} {
	return cmd.params
}

func (cmd *GetBackgroundColorsCommand) Done(result []byte, err error) {
	if cmd.cb == nil {
		return
	}
	if err != nil {
		cmd.cb(nil, err)
	} else {
		var rj GetBackgroundColorsResult
		if err := json.Unmarshal(result, &rj); err != nil {
			cmd.cb(nil, err)
		} else {
			cmd.cb(&rj, nil)
		}
	}
}

type MediaQueryResultChangedEvent struct {
}

// Fires whenever a MediaQuery result changes (for example, after a browser window has been resized.) The current implementation considers only viewport-dependent media features.
type MediaQueryResultChangedEventSink struct {
	events chan *MediaQueryResultChangedEvent
}

func NewMediaQueryResultChangedEventSink(bufSize int) *MediaQueryResultChangedEventSink {
	return &MediaQueryResultChangedEventSink{
		events: make(chan *MediaQueryResultChangedEvent, bufSize),
	}
}

func (s *MediaQueryResultChangedEventSink) Name() string {
	return "CSS.mediaQueryResultChanged"
}

func (s *MediaQueryResultChangedEventSink) OnEvent(params []byte) {
	evt := &MediaQueryResultChangedEvent{}
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

type FontsUpdatedEvent struct {
}

// Fires whenever a web font gets loaded.
type FontsUpdatedEventSink struct {
	events chan *FontsUpdatedEvent
}

func NewFontsUpdatedEventSink(bufSize int) *FontsUpdatedEventSink {
	return &FontsUpdatedEventSink{
		events: make(chan *FontsUpdatedEvent, bufSize),
	}
}

func (s *FontsUpdatedEventSink) Name() string {
	return "CSS.fontsUpdated"
}

func (s *FontsUpdatedEventSink) OnEvent(params []byte) {
	evt := &FontsUpdatedEvent{}
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

type StyleSheetChangedEvent struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"`
}

// Fired whenever a stylesheet is changed as a result of the client operation.
type StyleSheetChangedEventSink struct {
	events chan *StyleSheetChangedEvent
}

func NewStyleSheetChangedEventSink(bufSize int) *StyleSheetChangedEventSink {
	return &StyleSheetChangedEventSink{
		events: make(chan *StyleSheetChangedEvent, bufSize),
	}
}

func (s *StyleSheetChangedEventSink) Name() string {
	return "CSS.styleSheetChanged"
}

func (s *StyleSheetChangedEventSink) OnEvent(params []byte) {
	evt := &StyleSheetChangedEvent{}
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

type StyleSheetAddedEvent struct {
	Header *CSSStyleSheetHeader `json:"header"` // Added stylesheet metainfo.
}

// Fired whenever an active document stylesheet is added.
type StyleSheetAddedEventSink struct {
	events chan *StyleSheetAddedEvent
}

func NewStyleSheetAddedEventSink(bufSize int) *StyleSheetAddedEventSink {
	return &StyleSheetAddedEventSink{
		events: make(chan *StyleSheetAddedEvent, bufSize),
	}
}

func (s *StyleSheetAddedEventSink) Name() string {
	return "CSS.styleSheetAdded"
}

func (s *StyleSheetAddedEventSink) OnEvent(params []byte) {
	evt := &StyleSheetAddedEvent{}
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

type StyleSheetRemovedEvent struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // Identifier of the removed stylesheet.
}

// Fired whenever an active document stylesheet is removed.
type StyleSheetRemovedEventSink struct {
	events chan *StyleSheetRemovedEvent
}

func NewStyleSheetRemovedEventSink(bufSize int) *StyleSheetRemovedEventSink {
	return &StyleSheetRemovedEventSink{
		events: make(chan *StyleSheetRemovedEvent, bufSize),
	}
}

func (s *StyleSheetRemovedEventSink) Name() string {
	return "CSS.styleSheetRemoved"
}

func (s *StyleSheetRemovedEventSink) OnEvent(params []byte) {
	evt := &StyleSheetRemovedEvent{}
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

type LayoutEditorChangeEvent struct {
	StyleSheetId StyleSheetId `json:"styleSheetId"` // Identifier of the stylesheet where the modification occurred.
	ChangeRange  *SourceRange `json:"changeRange"`  // Range where the modification occurred.
}

type LayoutEditorChangeEventSink struct {
	events chan *LayoutEditorChangeEvent
}

func NewLayoutEditorChangeEventSink(bufSize int) *LayoutEditorChangeEventSink {
	return &LayoutEditorChangeEventSink{
		events: make(chan *LayoutEditorChangeEvent, bufSize),
	}
}

func (s *LayoutEditorChangeEventSink) Name() string {
	return "CSS.layoutEditorChange"
}

func (s *LayoutEditorChangeEventSink) OnEvent(params []byte) {
	evt := &LayoutEditorChangeEvent{}
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
