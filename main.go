package main

import (
	"log"
	"sync"

	"github.com/lkjfrf/content"
	"github.com/lkjfrf/core"
)

func main() {
	//core.GetLogManager().SetLogFile()
	content.GetContentManager().Init()
	core.GetNetworkCore().Init(":8002")
	log.Println("ScreenShareServer Start")

	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
}
