package setting

import (
	"sync"
)

type SettingHandler struct {
	ServerType int // 0: 나, 1: 원효로1번서버, 2: 원효로2번서버
	LogPath    string
}

var St_Ins *SettingHandler
var St_once sync.Once

func GetStManager() *SettingHandler {
	St_once.Do(func() {
		St_Ins = &SettingHandler{}
	})
	return St_Ins
}

func (st *SettingHandler) Init() {
	st.ServerType = 0 // 0: 나, 1: 원효로1번서버, 2: 원효로2번서버

	switch st.ServerType {
	case 0:
		st.LogPath = "/data/DIPServerLog/ScreenServer/"
	case 1:
		st.LogPath = "/data/DIPServerLog/ScreenServer1/"
	case 2:
		st.LogPath = "/data/DIPServerLog/ScreenServer2/"
	}
}
