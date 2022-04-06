package parser

type CSSDefinition struct {
	Selector *CSSSelector     `json:"selector"`
	Rules    []*CSSRule       `json:"rules"`
	Controls []*CSSDefinition `json:"controls"`
	DefLine  int              `json:"line"`
	Point    int              `json:"column"`
	Parent   *CSSDefinition   `json:"-"`
}

func NewDefinition(selector *CSSSelector, line, point int) *CSSDefinition { // 一个css定义
	return &CSSDefinition{
		Selector: selector,
		DefLine:  line,
		Point:    point - selector.RawOffset,
		Rules:    make([]*CSSRule, 0),
		Controls: make([]*CSSDefinition, 0),
	}
}

func (d *CSSDefinition) AddRule(rule *CSSRule) {
	d.Rules = append(d.Rules, rule)
}

func (d *CSSDefinition) AddControl(control *CSSDefinition) {
	d.Controls = append(d.Controls, control)
}

func (d *CSSDefinition) AddChild(def *CSSDefinition) {
	def.Parent = d
}

func (d *CSSDefinition) GetParent() *CSSDefinition {
	return d.Parent
}

func (d *CSSDefinition) IsControl() bool {
	return d.Selector.IsControlSelector()
}

type CSSDefinitionTree struct {
	definitions []*CSSDefinition
}

func NewDefinitionTree() *CSSDefinitionTree {
	return &CSSDefinitionTree{
		definitions: make([]*CSSDefinition, 0),
	}
}

func (l *CSSDefinitionTree) GetLastChild() *CSSDefinition {
	return l.definitions[len(l.definitions)-1]
}

func (l *CSSDefinitionTree) AddDefinitionToChild(def *CSSDefinition) {
	l.GetLastChild().AddChild(def)
}

func (l *CSSDefinitionTree) AddDefinition(def *CSSDefinition) {
	l.definitions = append(l.definitions, def)
}

func (l *CSSDefinitionTree) Remains() (remains bool) { // 是否还有剩余
	if len(l.definitions) > 0 {
		remains = true
	}
	return
}

func (l *CSSDefinitionTree) HasParent() (has bool) {
	if len(l.definitions) > 1 {
		has = true
	}
	return
}

func (l *CSSDefinitionTree) Remove() {
	l.definitions = l.definitions[0 : len(l.definitions)-1]
}

type CSSDefinitionList struct {
	definitions []*CSSDefinition
}

func NewDefinitionList() *CSSDefinitionList {
	return &CSSDefinitionList{
		definitions: make([]*CSSDefinition, 0),
	}
}

func (l *CSSDefinitionList) Add(def *CSSDefinition) {
	l.definitions = append(l.definitions, def)
}

func (l *CSSDefinitionList) Merge(defs []*CSSDefinition) {
	l.definitions = append(l.definitions, defs...)
}
func (l *CSSDefinitionList) Get() []*CSSDefinition {
	return l.definitions
}
