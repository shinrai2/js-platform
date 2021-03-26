package eventvm

/**
封装事件功能的js虚拟机原件
 */

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"time"
)

/**
 时钟原件
 */
type _timer struct {
	timer    *time.Timer
	duration time.Duration
	interval bool
	call     otto.FunctionCall
}
/**
 带事件的Javascript虚拟机
 */
type eventVM struct {
	ottoVm *otto.Otto
	registry map[*_timer]*_timer
	ready chan *_timer
}

func (vm *eventVM)VM() *otto.Otto {
	return vm.ottoVm
}

func (vm *eventVM)Dead() {
	close(vm.ready)
}

func newEventVM() *eventVM {
	event := &eventVM{
		ottoVm: otto.New(),
		registry: map[*_timer]*_timer{},
		ready: make(chan *_timer),
	}

	go func() {
		for timer := range event.ready {
			var arguments []interface{}
			if len(timer.call.ArgumentList) > 2 {
				tmp := timer.call.ArgumentList[2:]
				arguments = make([]interface{}, 2+len(tmp))
				for i, value := range tmp {
					arguments[i+2] = value
				}
			} else {
				arguments = make([]interface{}, 1)
			}
			arguments[0] = timer.call.ArgumentList[0]
			_, err := event.ottoVm.Call(`Function.call.call`, nil, arguments...)
			if err != nil {
				for _, timer := range event.registry {
					timer.timer.Stop()
					delete(event.registry, timer)
					fmt.Println(err)
				}
			}
			if timer.interval {
				timer.timer.Reset(timer.duration)
			} else {
				delete(event.registry, timer)
			}
		}
		fmt.Println("Event loop dead.")
	}()

	newTimer := func(call otto.FunctionCall, interval bool) (*_timer, otto.Value) {
		delay, _ := call.Argument(1).ToInteger()
		if 0 >= delay {
			delay = 1
		}

		timer := &_timer{
			duration: time.Duration(delay) * time.Millisecond,
			call:     call,
			interval: interval,
		}
		event.registry[timer] = timer

		timer.timer = time.AfterFunc(timer.duration, func() {
			event.ready <- timer
		})

		value, err := call.Otto.ToValue(timer)
		if err != nil {
			panic(err)
		}

		return timer, value
	}

	setTimeout := func(call otto.FunctionCall) otto.Value {
		_, value := newTimer(call, false)
		return value
	}
	event.ottoVm.Set("setTimeout", setTimeout)

	setInterval := func(call otto.FunctionCall) otto.Value {
		_, value := newTimer(call, true)
		return value
	}
	event.ottoVm.Set("setInterval", setInterval)

	clearTimeout := func(call otto.FunctionCall) otto.Value {
		timer, _ := call.Argument(0).Export()
		if timer, ok := timer.(*_timer); ok {
			timer.timer.Stop()
			delete(event.registry, timer)
		}
		return otto.UndefinedValue()
	}
	event.ottoVm.Set("clearTimeout", clearTimeout)
	event.ottoVm.Set("clearInterval", clearTimeout)

	return event
}
