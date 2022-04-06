package parser

import (
	"akt-wxss/utils"
)

type CSSRule struct {
	Property  string    `json:"property"` // 属性
	Value     *CSSValue `json:"value"`    // 值
	DefLine   int       `json:"line"`     // 行
	Point     int       `json:"column"`   // 位置
	RawData   string    `json:"raw"`      // 原始位置
	RawPoint  int       `json:"-"`
	RawOffset int       `json:"-"`
}

func NewRule(property []byte, line, point int) *CSSRule { // 属性
	_, prop, _, offset := utils.ParseBytes(property)
	return &CSSRule{
		Property:  string(prop),
		DefLine:   line,
		Point:     point - offset,
		RawData:   string(property),
		RawPoint:  point,
		RawOffset: offset,
	}
}

func (r *CSSRule) IsSpecialProperty() (special bool) {
	if r.Property == "filter" {
		special = true
	}
	// todo
	return
}

func (r *CSSRule) SetValue(value []byte, index, point int, semicolon bool) {
	r.Value = NewValue(value, index, point, semicolon)
}
