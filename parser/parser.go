package parser

import (
	"akt-wxss/utils"
	"bytes"
)

type CSSValue struct {
	Value     string `json:"value"`     // 值
	DefLine   int    `json:"line"`      // 行
	Point     int    `json:"column"`    // 位置
	Semicolon bool   `json:"semicolon"` // 是否有分号
	RawData   string `json:"raw"`       // 原始数据
}

func NewValue(val []byte, line, point int, semicolon bool) *CSSValue {
	_, value, _, offset := utils.ParseBytes(val)
	return &CSSValue{
		Value:     string(value),
		DefLine:   line,
		Point:     point - offset,
		RawData:   string(val),
		Semicolon: semicolon,
	}
}

type CSSParseResult struct {
	data []*CSSDefinition
}

func (c CSSParseResult) GetData() []*CSSDefinition {
	return c.data
}

func (c CSSParseResult) Traverse(callbacks ...func(*CSSDefinition)) {
	callback := func(node *CSSDefinition) {
		for _, cbb := range callbacks {
			cbb(node)
		}
	}
	for _, v := range c.data {
		callback(v)
		if len(v.Controls) > 0 {
			for _, co := range v.Controls {
				callback(co)
			}
		}
	}

}

func (c CSSParseResult) Minisize() bytes.Buffer { // 压缩
	var buffer bytes.Buffer
	for _, v := range c.data {
		s := v.Selector.Selector
		buffer.WriteString(s)
		buffer.WriteByte(utils.BRACE_OPEN)
		for _, r := range v.Rules {
			p := r.Property
			v := r.Value.Value
			buffer.WriteString(p)
			buffer.WriteByte(utils.COLON)
			buffer.WriteString(v)
			buffer.WriteByte(utils.SEMI)
		}
		buffer.WriteByte(utils.BRACE_CLOSE)

	}
	return buffer
}

type CSSParser struct {
	definitions *CSSDefinitionList
	defTree     *CSSDefinitionTree
	defRule     *CSSRule
	stack       []byte

	charPoint   int
	comment     bool
	quoting     bool
	singleQuote bool
	doubleQuote bool
	inSelector  bool
	skipping    bool
	isEscaping  bool
	inParen     bool
}

func NewParser() *CSSParser {
	return &CSSParser{
		definitions: NewDefinitionList(),
		stack:       []byte{},
		defTree:     NewDefinitionTree(),
		defRule:     nil,
		charPoint:   0,
	}
}

// 注释
func (c *CSSParser) isCommentStart(line []byte, point int) bool {
	if len(line) <= point+1 || c.quoting {
		return false
	}

	if line[point] == utils.COMMENT_SLASH && line[point+1] == utils.COMMENT_STAR {
		return true
	}

	return false
}

func (c *CSSParser) isCommentEnd(line []byte, point int) bool {
	if len(line) <= point+1 || c.quoting {
		return false
	}

	if line[point] == utils.COMMENT_STAR && line[point+1] == utils.COMMENT_SLASH {
		return true
	}

	return false
}

// 转义
func (c *CSSParser) parseEscapeSequence() {
	if !c.skipping {
		c.skipping = true
		c.isEscaping = true
	} else {
		c.skipping = false
		c.isEscaping = false
	}

	c.stack = append(c.stack, utils.ESCAPE_SEQUENCE)
}

// 换行
func (c *CSSParser) parseLineFeed(index *int) {
	val := bytes.Trim(c.stack, ";:\n\t ")

	if len(val) > 0 {
		if !c.inSelector && val[0] == utils.AT {
			def := NewDefinition(
				NewSelector(c.stack),
				*index,
				c.charPoint,
			)
			c.definitions.Add(def)
			c.stack = []byte{}
		} else if c.defRule != nil { // 设置属性
			c.defRule.SetValue(
				c.stack,
				*index,
				c.charPoint,
				false,
			)
			c.defTree.GetLastChild().AddRule(c.defRule)
			c.defRule = nil
			c.stack = []byte{}
		}
	}
	c.charPoint = 0
}

// 匹配{ 进入选择器
func (c *CSSParser) parseBraceOpen(index *int) {
	def := NewDefinition(
		NewSelector(c.stack),
		*index,
		c.charPoint,
	)
	c.defTree.AddDefinition(def)
	c.inSelector = true
	c.stack = []byte{}
}

// :
func (c *CSSParser) parseColon(index *int) {
	if c.defRule != nil && c.defRule.IsSpecialProperty() || !c.inSelector {
		c.stack = append(c.stack, utils.AT)
		return
	}
	c.defRule = NewRule(c.stack, *index, c.charPoint)
	c.stack = []byte{}
}

func isEmptyStack(stack []byte) (isEmpty bool) {
	if len(bytes.Trim(stack, "\r\n\t ")) == 0 {
		isEmpty = true
	}
	return
}

// ;
func (c *CSSParser) parseSemi(index *int) {
	if !c.inSelector {
		def := NewDefinition(
			NewSelector(c.stack),
			*index,
			c.charPoint,
		)
		c.definitions.Add(def)
		c.stack = []byte{}
		return
	}
	if !isEmptyStack(c.stack) {
		c.defRule.SetValue(
			c.stack,
			*index,
			c.charPoint,
			true,
		)
		c.defTree.GetLastChild().AddRule(c.defRule)
		c.defRule = nil
	}
	c.stack = []byte{}
}

// }
func (c *CSSParser) parseBraceClose(index *int) {
	cdef := c.defTree.GetLastChild()
	if c.defRule != nil {
		c.defRule.SetValue(
			c.stack,
			*index,
			c.charPoint,
			false,
		)
		cdef.AddRule(c.defRule)
		c.defRule = nil
	}
	c.defTree.Remove()

	if c.defTree.Remains() {
		c.defTree.GetLastChild().AddControl(cdef)
	} else {
		c.definitions.Add(cdef)
	}
	c.inSelector = false
	c.stack = []byte{}
}

func (c *CSSParser) execParse(line []byte) {
	index := 1

	for point := 0; point < len(line); point++ {
		c.charPoint++

		if c.isCommentStart(line, point) {
			c.comment = true
			c.stack = append(c.stack, line[point])
			continue
		}

		if c.isCommentEnd(line, point) {
			c.comment = false
			c.stack = append(c.stack, utils.COMMENT_STAR, utils.COMMENT_SLASH)
			point++
			c.charPoint++
			continue
		}

		if c.comment {
			c.stack = append(c.stack, line[point])
			continue
		}

		switch line[point] {
		case utils.ESCAPE_SEQUENCE:
			c.parseEscapeSequence()
			continue
		case utils.PAREN_LEFT:
			if !c.quoting {
				c.inParen = true
			}
		case utils.PAREN_RIGHT:
			if !c.quoting {
				c.inParen = false
			}
		case utils.LINE_FEED:
			c.parseLineFeed(&index)
			index++
		case utils.DOUBLE_QUOTE:
			if c.skipping || c.singleQuote {
				break
			}
			if c.doubleQuote {
				c.quoting = false
				c.doubleQuote = false
			} else {
				c.quoting = true
				c.doubleQuote = true
			}
		case utils.SINGLE_QUOTE:
			if c.skipping || c.doubleQuote {
				break
			}
			if c.singleQuote {
				c.quoting = false
				c.singleQuote = false
			} else {
				c.quoting = true
				c.singleQuote = true
			}
		case utils.BRACE_OPEN:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseBraceOpen(&index)
			continue
		case utils.COLON:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseColon(&index)
			continue
		case utils.SEMI:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseSemi(&index)
			continue
		case utils.BRACE_CLOSE:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseBraceClose(&index)
			continue
		}

		if c.isEscaping {
			c.isEscaping = false
			c.skipping = false
		}

		c.stack = append(c.stack, line[point])

	}

	if c.defRule != nil {
		c.defRule.SetValue(c.stack, index, c.charPoint, false)
		c.defTree.GetLastChild().AddRule(c.defRule)
		c.defRule = nil
	}

	if c.defTree.Remains() {
		c.definitions.Add(c.defTree.GetLastChild())
	}

}

func (c *CSSParser) Parse(buffer []byte) CSSParseResult {
	c.execParse(buffer)

	return CSSParseResult{
		data: c.definitions.Get(),
	}
}
