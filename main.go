package main

import (
	"github.com/lkjfrf/content"
	"github.com/lkjfrf/core"
)

func main() {
	core.GetLogManager().SetLogFile()
	core.GetNetworkCore().Init(":8005")
	content.GetContentManager().Init()
	for {

	}
}
