package main

import (
	"github.com/jeffchannell/lbdl/main/lbdl"
)

func main() {
	if err := lbdl.Start(); nil != err {
		panic(err)
	}
}
