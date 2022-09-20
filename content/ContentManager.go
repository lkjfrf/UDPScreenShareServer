package content

import (
	"encoding/json"

	"net"
	"sync"
)

type ContentManager struct {
	HandlerFunc    map[int]func(*net.UDPConn, *net.UDPAddr, string)
	ScreenChannels sync.Map // 스크린 공유 채널목록들 안에 ScreenShares 저장
	ScreenShares   sync.Map // 스크린 공유자들 안에 스크린 관람자들 ID 저장
}

var CM_Ins *ContentManager
var CM_once sync.Once

func GetContentManager() *ContentManager {
	CM_once.Do(func() {
		CM_Ins = &ContentManager{}
	})
	return CM_Ins
}

func (cm *ContentManager) Init() {
	cm.HandlerFunc = make(map[int]func(*net.UDPConn, *net.UDPAddr, string), 0)
	cm.ScreenChannels = sync.Map{}
	cm.ScreenShares = sync.Map{}

	cm.HandlerFunc[ChannelEnter] = cm.ChannelEnter
	cm.HandlerFunc[Voice] = cm.Voice
	cm.HandlerFunc[EScreenShareToggle] = cm.ScreenShareToggle
	cm.HandlerFunc[EScreenShareView] = cm.ScreenShareView
	cm.HandlerFunc[PlayerLogout] = cm.PlayerLogout

	cm.Test()
}

func (cm *ContentManager) ChannelEnter(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ChannelEnter{}
	json.Unmarshal([]byte(jsonstr), &data)

	if data.ChannelType == 2 {
		cm.ScreenChannels.LoadOrStore(data.ChannelNum, &sync.Map{})
	}
	GetSession().NewPlayer(data.Id, conn, addr, data.ChannelNum)
}

func (cm *ContentManager) Voice(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := SR_Voice{}
	json.Unmarshal([]byte(jsonstr), &data)

	GetSession().BroadCastToSameChannelExpetMe(data.Id, data, Voice)
}

func (cm *ContentManager) ScreenShareToggle(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ScreenShareToggle{}
	json.Unmarshal([]byte(jsonstr), &data)

	if Ch, ok := cm.ScreenChannels.Load(data.ChannelNum); ok {
		if data.IsOn {
			Ch.(*sync.Map).LoadOrStore(GetSession().GetPlayerIdByCon(conn), &sync.Map{})

		}
	}
	//GetSession().BroadCastToSameChannelExpetMe(data.Id, data, Voice)
}
func (cm *ContentManager) ScreenShareView(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ScreenShareView{}
	json.Unmarshal([]byte(jsonstr), &data)

	//	GetSession().BroadCastToSameChannelExpetMe(data.Id, data, Voice)
}
func (cm *ContentManager) PlayerLogout(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_PlayerLogout{}
	json.Unmarshal([]byte(jsonstr), &data)

	GetSession().BroadCastToSameChannelExpetMe(data.Id, data, Voice)
}

func (cm *ContentManager) Test() {
	cm.ScreenChannels.LoadOrStore(1, &sync.Map{})
	cm.ScreenChannels.LoadOrStore(2, &sync.Map{})
	cm.ScreenChannels.LoadOrStore(1, &sync.Map{})

	if Ch, ok := cm.ScreenChannels.Load(1); ok {
		if true {
			if ChPlayer, ok := Ch.(*sync.Map).LoadOrStore("song", &sync.Map{}); ok {
				ChPlayer.(*sync.Map).Store("hihi", true)
			}
		}
	}
	if Ch, ok := cm.ScreenChannels.Load(1); ok {
		if true {
			if ChPlayer, ok := Ch.(*sync.Map).LoadOrStore("song", &sync.Map{}); ok {
				ChPlayer.(*sync.Map).Store("noo", false)
			}
		}
	}
}
