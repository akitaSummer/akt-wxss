package utils

import "regexp"

const (
	BRACE_OPEN      = 123 // "{"
	BRACE_CLOSE     = 125 // "}"
	COLON           = 58  // ":"
	SEMI            = 59  // ";"
	COMMENT_SLASH   = 47  // "/"
	COMMENT_STAR    = 42  // "*"
	AT              = 64  // "@"
	DOUBLE_QUOTE    = 34  // "\""
	SINGLE_QUOTE    = 39  // "'"
	PAREN_LEFT      = 40  // "("
	PAREN_RIGHT     = 41  // ")"
	LINE_FEED       = 10  // "\n"
	CARRIAGE_RETURN = 13  // "\r"
	ESCAPE_SEQUENCE = 92  // "\"
)

var (
	leftRegex    = regexp.MustCompile("^([;\n\t ]*)")
	rightRegex   = regexp.MustCompile("[\n\t/\\* ]*$")
	spaceRegex   = regexp.MustCompile("[\n ]+")
	commentRegex = regexp.MustCompile("/\\*.*\\*/")
)

func ParseBytes(data []byte) (before, value, after []byte, offset int) {
	size := len(data)
	data = commentRegex.ReplaceAll(data, []byte(""))
	left := leftRegex.FindSubmatchIndex(data)
	before = data[:left[1]]
	data = data[left[1]:]

	right := rightRegex.FindSubmatchIndex(data)
	value = data[:right[0]]
	after = data[right[0]:]

	offset = size - len(before)
	value = spaceRegex.ReplaceAll(value, []byte(" "))

	return
}
