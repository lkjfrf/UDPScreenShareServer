package content

import (
	"encoding/json"
	"log"

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
	cm.HandlerFunc[EScreenShare] = cm.ScreenShare
	cm.HandlerFunc[EScreenShareToggle] = cm.ScreenShareToggle
	cm.HandlerFunc[EScreenShareView] = cm.ScreenShareView
	cm.HandlerFunc[PlayerLogout] = cm.PlayerLogout

	cm.Test()
}

func (cm *ContentManager) ChannelEnter(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ChannelEnter{}
	json.Unmarshal([]byte(jsonstr), &data)

	// if data.ChannelType == 2 {
	// 	cm.ScreenChannels.LoadOrStore(data.ChannelNum, &sync.Map{})
	// }
	GetSession().NewPlayer(data.Id, conn, addr, data.ChannelNum)

	// _, AlreadyLogin := GetSession().GSession.Load(data.Id)
	// if !AlreadyLogin {
	// 	GetSession().NewPlayer(data.Id, conn, addr, data.ChannelNum)
	// } else {
	// 	if p, ok := GetSession().GSession.Load(data.Id); ok {
	// 		p.(*Player).Channel = data.ChannelNum
	// 	}
	// }

}

func (cm *ContentManager) ScreenShare(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {

	// id := GetSession().GetPlayerIdByCon(conn)
	// channelNum := instance_gs.GetChannelNumById(id)
	// if Ch, ok := cm.ScreenChannels.Load(channelNum); ok {
	// 	if ChSharer, ok := Ch.(*sync.Map).Load(id); ok {
	// 		ChSharer.(*sync.Map).Range(func(key, value any) bool {

	// 			GetSession().SendByte(conn, addr, []byte(jsonstr))
	// 			return true
	// 		})
	// 	}
	// }
	data := SR_ScreenShare{}
	json.Unmarshal([]byte(jsonstr), &data)

	GetSession().BroadCastToSameChannelNum(GetSession().GetChannelNumById(data.Id), data, EScreenShare)
}

func (cm *ContentManager) ScreenShareToggle(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ScreenShareToggle{}
	json.Unmarshal([]byte(jsonstr), &data)
	myId := GetSession().GetPlayerIdByCon(conn)
	ChannelNum := GetSession().GetChannelNumById(myId)

	if Ch, ok := cm.ScreenChannels.Load(ChannelNum); ok {
		if data.IsOn {
			Ch.(*sync.Map).LoadOrStore(myId, &sync.Map{})
		} else {
			Ch.(*sync.Map).Delete(myId)
		}
	}
	log.Println(myId, "Has ", data.IsOn, " ShareScreen")
}

func (cm *ContentManager) ScreenShareView(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ScreenShareView{}
	json.Unmarshal([]byte(jsonstr), &data)
	myId := GetSession().GetPlayerIdByCon(conn)
	if Ch, ok := cm.ScreenChannels.Load(GetSession().GetChannelNumById(myId)); ok {
		if ChSharer, ok := Ch.(*sync.Map).Load(data.ViewTarget); ok {
			if data.IsOn {
				ChSharer.(*sync.Map).Store(myId, myId)

				packet := R_ScreenShareView{IsHasViewer: true}
				GetSession().SendPacketById(myId, packet, EScreenShareView)
			} else {
				ChSharer.(*sync.Map).Delete(myId)
			}
		}
	}
}

func (cm *ContentManager) PlayerLogout(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_PlayerLogout{}
	json.Unmarshal([]byte(jsonstr), &data)
	ChannelNum := GetSession().GetChannelNumById(data.Id)

	if Ch, ok := cm.ScreenChannels.Load(ChannelNum); ok {
		Ch.(*sync.Map).Delete(data.Id)
	}
	GetSession().GSession.Delete(data.Id)
	GetSession().ConMap.Delete(conn)
	log.Println(data.Id, " Log out")
}

func (cm *ContentManager) Test() {
	// a := []uint8{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 255, 219, 0, 67, 0, 80, 55, 60, 70, 60, 50, 80, 70, 65, 70, 90, 85, 80, 95, 120, 200, 130, 120, 110, 110, 120, 245, 175, 185, 145, 200, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	// log.Println(a)
	// b := []uint16{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 255, 219, 0, 67, 0, 80, 55, 60, 70, 60, 50, 80, 70, 65, 70, 90, 85, 80, 95, 120, 200, 130, 120, 110, 110, 120, 245, 175, 185, 145, 200, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	// log.Println(b)

	// jsona := "{\"id\":\"q1\",\"data\":[255,216,255,224,0,16,74,70,73,70]}"
	// log.Println(jsona)
	// type st8 struct {
	// 	Data []uint8
	// }
	// type st16 struct {
	// 	Data []uint16
	// }

	// data := st8{}
	// data2 := st16{}
	// json.Unmarshal([]byte(jsona), &data)
	// json.Unmarshal([]byte(jsona), &data2)
	// log.Println(data.Data)
	// log.Println(data2.Data)
}
