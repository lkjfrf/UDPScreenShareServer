package main

import (
	"time"

	"github.com/lkjfrf/content"
	"github.com/lkjfrf/core"
)

func main() {
	//core.GetLogManager().SetLogFile()
	content.GetContentManager().Init()
	core.GetNetworkCore().Init(":8005")
	time.Sleep(time.Minute * 100)
}
