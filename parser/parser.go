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
