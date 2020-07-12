// +build windows

package main

import (
	"os"
	"strings"
)

func programPrefix() string {
	p := os.Args[0]
	if strings.HasSuffix(strings.ToLower(p), ".exe") {
		p = p[len(p)-5:]
	}
	return p
}
