package main

/**
交互控制
 */

import (
	"bufio"
	"os"
)

type IOHandle struct {
	regOutFunc func()
	regInFunc []func(string)bool
}

func (handle *IOHandle)channelHandleOut(fn func()) {
	handle.regOutFunc = fn
}

func (handle *IOHandle)channelHandleIn(startWith string, fn func(string)) {
	handle.regInFunc = append(handle.regInFunc, func(s string)bool {
		sw := startWith
		fn1 := fn
		if len(s) >= len(sw) && s[:len(sw)] == sw {
			fn1(s[len(sw):])
			return true
		}
		return false
	})
}

func (handle *IOHandle)loop() {
	go handle.regOutFunc()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		for _, fn := range handle.regInFunc {
			if fn(s) {
				break
			}
		}
	}
}

