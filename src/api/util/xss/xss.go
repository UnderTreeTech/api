package xss

import (
	"github.com/microcosm-cc/bluemonday"
)

var p *bluemonday.Policy

func init() {
	p = bluemonday.UGCPolicy()
}

func GetXssHandler() *bluemonday.Policy {
	return p
}
