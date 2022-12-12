package main

import (
	"log"
	"sync"

	"github.com/ScreenShare/content"
	"github.com/ScreenShare/core"
	"github.com/ScreenShare/setting"
)

func main() {
	setting.GetStManager().Init()
	setting.GetLogManager().SetLogFile()
	content.GetContentManager().Init()
	core.GetNetworkCore().Init(setting.St_Ins.Port)
	log.Println("ScreenShareServer Start")

	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
}
