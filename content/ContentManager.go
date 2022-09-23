package content

import (
	"encoding/json"
	"log"

	"net"
	"sync"
)

type ContentManager struct {
	HandlerFunc map[int]func(*net.UDPConn, *net.UDPAddr, string)
}

type Connection struct {
	Con  *net.UDPConn
	Addr *net.UDPAddr
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

	cm.HandlerFunc[ChannelEnter] = cm.ChannelEnter
	cm.HandlerFunc[EScreenShare] = cm.ScreenShare
	cm.HandlerFunc[PlayerLogout] = cm.PlayerLogout
	cm.HandlerFunc[EScreenWatchToggle] = cm.ScreenWatchToggle

	cm.Test()
}

func (cm *ContentManager) ChannelEnter(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ChannelEnter{}
	json.Unmarshal([]byte(jsonstr), &data)

	Con := Connection{Con: conn, Addr: addr}
	_, AlreadyLogin := GetSession().ConMap.Load(Con)

	if !AlreadyLogin {
		GetSession().NewPlayer(data.Id, Con, data.ChannelNum)
	} else {
		if p, ok := GetSession().Players.Load(data.Id); ok {
			p.(*Player).Channel = data.ChannelNum
		}
	}
}

func (cm *ContentManager) ScreenShare(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := SR_ScreenShare{}
	json.Unmarshal([]byte(jsonstr), &data)
	log.Println(data.Id, "- ", data.Sequence)
	sendBuffer := MakeSendBuffer(EScreenShare, data)

	//	var SendQueue []Connection
	targetChannel := GetSession().GetChannelNumById(data.Id)
	GetSession().Players.Range(func(key, value any) bool {
		if value.(*Player).ScreenOn && value.(*Player).Channel == targetChannel {
			//if value.(*Player).Channel == GetSession().GetChannelNumById(data.Id) {
			GetSession().SendByte(value.(*Player).Conn, value.(*Player).Addr, sendBuffer)
			//SendQueue = append(SendQueue, Connection{Con: value.(*Player).Conn, Addr: value.(*Player).Addr})

		}
		return true
	})

	// maxHandle := 100
	// totalSendLen := len(SendQueue)
	// if totalSendLen == 0 {
	// 	return
	// }
	// for j := 0; j <= totalSendLen/maxHandle; j++ {
	// 	go func(j int) { // 100개의 패킷 전송당 1개의 go 루틴
	// 		log.Println(j)
	// 		if totalSendLen/maxHandle > j {
	// 			for i := j * maxHandle; i < (j+1)*maxHandle; i++ {
	// 				GetSession().SendByte(SendQueue[i].Con, SendQueue[i].Addr, sendBuffer)
	// 			}
	// 		} else if totalSendLen/maxHandle == j { // 마지막 고루틴
	// 			for i := j * maxHandle; i < (j*maxHandle)+(totalSendLen-j*maxHandle); i++ {
	// 				GetSession().SendByte(SendQueue[i].Con, SendQueue[i].Addr, sendBuffer)
	// 			}
	// 		}
	// 	}(j)
	// }

	// for i := 0; i < len(SendQueue); i++ {
	// 	GetSession().SendByte(SendQueue[i].Con, SendQueue[i].Addr, sendBuffer)
	// }

	//GetSession().BroadCastToSameChannelNum(GetSession().GetChannelNumById(data.Id), data, EScreenShare)
}

func (cm *ContentManager) ScreenWatchToggle(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_ScreenWatchToggle{}
	json.Unmarshal([]byte(jsonstr), &data)
	//id := GetSession().GetPlayerIdByCon(Connection{Con: conn, Addr: addr})
	if p, ok := GetSession().Players.Load(data.Id); ok {
		p.(*Player).ScreenOn = data.IsOn
	}
}

func (cm *ContentManager) PlayerLogout(conn *net.UDPConn, addr *net.UDPAddr, jsonstr string) {
	data := S_PlayerLogout{}
	json.Unmarshal([]byte(jsonstr), &data)

	GetSession().Players.Delete(data.Id)
	GetSession().ConMap.Delete(Connection{Con: conn, Addr: addr})
	log.Println(data.Id, " Log out")
}

func (cm *ContentManager) Test() {
	// GetSession().Players.Store("q", 123)
	// if p, ok := GetSession().Players.LoadAndDelete("q"); ok {
	// 	log.Println(p)
	// }
	// if p, ok := GetSession().Players.Load("q"); ok {
	// 	log.Println(p)
	// }

}
