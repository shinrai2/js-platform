package eventvm

/**
js虚拟机封装和管理器
 */

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"sync"
	"time"
)

var controlLock sync.Mutex

const NotSelected int64 = -1

type vmMaster struct {
	vms map[int64]*unixVM
	current int64
	inCh chan string
	outCh chan string
}

type unixVM struct {
	unix int64
	inCh chan string
	outCh chan string
}

func vmCreate() (chan string, chan string) {
	inCh := make(chan string)
	outCh := make(chan string)
	go func() {
		defer close(outCh)
		vm := newEventVM()
		defer vm.Dead()
		// 平台方法
		vm.VM().Set("console", map[string]interface{}{
			"log": func(c otto.FunctionCall) otto.Value {
				for _, argsOfC := range c.ArgumentList {
					outCh <- fmt.Sprintf("[%s] %s",
						time.Now().Format("2006-01-02 15:04:05"), argsOfC.String())
				}
				return otto.UndefinedValue()
			},
		})
		vm.VM().Set("platform", map[string]interface{}{
			"log": func(c otto.FunctionCall) otto.Value {
				for _, argsOfC := range c.ArgumentList {
					fmt.Printf("[%s] %s\n",
						time.Now().Format("2006-01-02 15:04:05"), argsOfC.String())
				}
				return otto.UndefinedValue()
			},
		})
		for in := range inCh {
			vm.VM().Run(in)
		}
	}()
	return inCh, outCh
}

func NewMaster() *vmMaster {
	master := vmMaster{
		vms:     map[int64]*unixVM{},
		current: NotSelected,
		inCh:    make(chan string),
		outCh:   make(chan string),
	}
	go func() {
		for in := range master.inCh {
			if master.current == NotSelected {
				fmt.Println("Undefined operation.")
				// DO NOTHING.
			} else {
				master.vms[master.current].inCh <- in
			}
		}
		for k := range master.vms {
			close(master.vms[k].inCh) // 通知关闭所有vm
		}
	}()
	return &master
}

func (master *vmMaster)controlVM(inCh chan string, outCh chan string) int64 {
	controlLock.Lock()
	now := time.Now().Unix()
	controlLock.Unlock()
	pack := unixVM{
		unix:  now,
		inCh:  inCh,
		outCh: outCh,
	}
	master.vms[now] = &pack
	go func(u int64) {
		for out := range outCh {
			if master.current == u {
				master.outCh <- out
			} else {
				// DO NOTHING.
			}
		}
	}(now)
	return now
}

func (master *vmMaster)CreateVM() int64 {
	in, out := vmCreate()
	return master.controlVM(in, out)
}

func (master *vmMaster)GetList() []int64 {
	keys := make([]int64, 0, len(master.vms))
	for k := range master.vms {
		keys = append(keys, k)
	}
	return keys
}

func (master *vmMaster)Switch(unix int64) {
	if _, ok := master.vms[unix]; ok || unix == NotSelected {
		master.current = unix
	}
}

func (master *vmMaster)GetIO() (chan string, chan string) {
	return master.inCh, master.outCh
}

func (master *vmMaster)Current() int64 {
	return master.current
}


