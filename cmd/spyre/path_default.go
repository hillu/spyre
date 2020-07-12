// +build !windows

package main

import (
	"os"
)

func programPrefix() string { return os.Args[0] }
