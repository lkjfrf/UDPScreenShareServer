package main

import (
	"time"

	"github.com/lkjfrf/content"
	"github.com/lkjfrf/core"
)

func main() {
	core.GetLogManager().SetLogFile()
	core.GetNetworkCore().Init(":8001")
	content.GetContentManager().Init()
	time.Sleep(time.Minute * 100)
}
