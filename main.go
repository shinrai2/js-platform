package main
/**
程序入口
 */

import (
	"fmt"
	"js-platform/eventvm"
	"os"
	"strconv"
)

func main() {

	master := eventvm.NewMaster()
	in, out := master.GetIO()

	ioHandle := IOHandle{}
	ioHandle.channelHandleOut(func() {
		fmt.Println("# Output start.")
		for o := range out {
			fmt.Println(o)
		}
		fmt.Println("# Output end.")
	})
	ioHandle.channelHandleIn("new.", func(_ string) {
		fmt.Println("ID:", master.CreateVM())
	})
	ioHandle.channelHandleIn("list.", func(_ string) {
		fmt.Println("List: ", master.GetList())
	})
	ioHandle.channelHandleIn("current.", func(_ string) {
		fmt.Println("Current: ", master.Current())
	})
	ioHandle.channelHandleIn("switch.", func(s string) {
		var i64 int64
		if s == "" { // switch to default.
			i64 = eventvm.NotSelected
		} else {
			i, _ := strconv.Atoi(s)
			i64 = int64(i)
		}
		master.Switch(i64)
		fmt.Println(master.Current() == i64)
	})
	ioHandle.channelHandleIn("exit.", func(_ string) {
		close(in)
		os.Exit(0)
	})
	ioHandle.channelHandleIn("", func(s string) {
		in <- s
	})

	ioHandle.loop()
}
