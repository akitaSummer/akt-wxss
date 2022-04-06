package parser

import (
	"akt-wxss/utils"
)

type CSSSelector struct {
	Selector        string `json:"selector"`
	ControlSelector bool   `json:"atrule"`
	RawData         string `json:"raw"`
	RawOffset       int    `json:"-"`
}

func NewSelector(selBytes []byte) *CSSSelector { // 选择器
	_, selector, _, offset := utils.ParseBytes(selBytes)
	var isControl bool

	if len(selector) > 0 && selector[0] == utils.AT { // @
		isControl = true
	} else {
		isControl = false
	}

	return &CSSSelector{
		Selector:        string(selector),
		ControlSelector: isControl,
		RawData:         string(selBytes),
		RawOffset:       offset,
	}
}

func (s *CSSSelector) String() string {
	return s.Selector
}

func (s *CSSSelector) IsControlSelector() bool {
	return s.ControlSelector
}
