package parser

import (
	"akt-wxss/utils"
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
