package provider

import (
	"bytes"
)

// Driver ...
type Driver interface {
	Send(mobile, content string) (string, error)
	Assignment(content string, value ...string) string
	String() string
	GetDebug() bool
}

// SignContent ...
func SignContent(sign, content string) string {
	var buff bytes.Buffer
	buff.WriteString("【")
	buff.WriteString(sign)
	buff.WriteString("】")
	buff.WriteString(content)
	return buff.String()
}
