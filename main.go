package main

import (
	"godaemon/core"
	"godaemon/initialization"
)

func main() {
	err := initialization.Init()
	if err != nil {
		panic(err)
	}
	go initialization.UpTask()
	core.RunServer()

}
