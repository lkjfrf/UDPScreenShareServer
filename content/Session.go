package content

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"sync"
)

type GlobalSession struct {
	Players sync.Map
	ConMap  sync.Map
	//Channels map[int32]Player
}

var instance_gs *GlobalSession
var once_gs sync.Once

func GetSession() *GlobalSession {
	once_gs.Do(func() {
		instance_gs = &GlobalSession{}
	})
	return instance_gs
}

type Player struct {
	Conn     *net.UDPConn
	Addr     *net.UDPAddr
	Channel  int32
	ScreenOn bool
}

func (gs *GlobalSession) Init() {
	gs.Players = sync.Map{}
	gs.ConMap = sync.Map{}
}

func (gs *GlobalSession) NewPlayer(id string, c Connection, channelNum int32) {
	Con_player := &Player{}
	Con_player.Conn = c.Con
	Con_player.Addr = c.Addr
	Con_player.Channel = channelNum
	Con_player.ScreenOn = false
	gs.Players.Store(id, Con_player)
	gs.ConMap.Store(c, id)
	log.Println(id, "Enter", channelNum)
}

func (gs *GlobalSession) GetPlayerIdByCon(c Connection) string {
	if id, ok := gs.ConMap.Load(c); ok {
		return id.(string)
	}
	return ""
}

func (gs *GlobalSession) GetChannelNumById(Id string) int32 {
	if ch, ok := gs.Players.Load(Id); ok {
		return ch.(*Player).Channel
	}
	return -1
}

func (gs *GlobalSession) BroadCastToSameChannelNum(ChannelNum int32, recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)

	gs.Players.Range(func(key, value any) bool {
		if value.(*Player).Channel == ChannelNum {
			gs.SendByte(value.(*Player).Conn, value.(*Player).Addr, sendBuffer)
			//gs.SendByte2(value.(*Player).Conn, data)
		}
		return true
	})
}

func (gs *GlobalSession) BroadCastToSameChannelNumExpetMe(ChannelNum int32, Id string, data []byte) {
	gs.Players.Range(func(key, value any) bool {
		if value.(*Player).Channel == ChannelNum && key != Id {
			gs.SendByte(value.(*Player).Conn, value.(*Player).Addr, data)
		}
		return true
	})
}

func (gs *GlobalSession) BroadCastToSameChannelExpetMe(id string, data []byte) {
	var TargetChannel int32
	if p, ok := gs.Players.Load(id); ok {
		TargetChannel = p.(*Player).Channel
	}
	gs.Players.Range(func(key, value any) bool {
		if value.(*Player).Channel == TargetChannel && key != id {
			gs.SendByte(value.(*Player).Conn, value.(*Player).Addr, data)
		}
		return true
	})
}

func (gs *GlobalSession) SendByte(c *net.UDPConn, addr *net.UDPAddr, data []byte) {
	if c != nil {
		sent, err := c.WriteToUDP(data, addr)
		if err != nil {
			log.Println("SendPacket ERROR :", err)
		} else {
			if sent != len(data) {
				log.Println("[Sent diffrent size] : SENT =", sent, "BufferSize =", len(data))
			}
			//log.Println("SendPacket ", addr, "/", c)
		}
	}
}
func (gs *GlobalSession) SendByte2(c *net.UDPConn, data []byte) {
	if c != nil {
		sent, err := c.Write(data)
		if err != nil {
			log.Println("SendPacket ERROR :", err)
		} else {
			if sent != len(data) {
				log.Println("[Sent diffrent size] : SENT =", sent, "BufferSize =", len(data))
			}
			//log.Println("SendPacket ", addr, "/", c)
		}
	}
}

func (gs *GlobalSession) SendPacketByConn2(conn *net.UDPConn, recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)

	gs.SendByte2(conn, sendBuffer)
}

func (gs *GlobalSession) SendPacketByConn(conn *net.UDPConn, addr *net.UDPAddr, recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)

	gs.SendByte(conn, addr, sendBuffer)
}

func (gs *GlobalSession) SendPacketById(Id string, recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)

	if p, ok := gs.Players.Load(Id); ok {
		gs.SendByte(p.(*Player).Conn, p.(*Player).Addr, sendBuffer)
	}
}

func MakeSendBuffer[T any](pktid uint16, data T) []byte {
	sendData, err := json.Marshal(&data)
	if err != nil {
		log.Println("MakeSendBuffer : Marshal Error", err)
	}
	sendBuffer := make([]byte, 4)

	pktsize := len(sendData) + 4

	binary.LittleEndian.PutUint32(sendBuffer, uint32(pktsize))
	binary.LittleEndian.PutUint16(sendBuffer[2:], pktid)

	sendBuffer = append(sendBuffer, sendData...)

	return sendBuffer
}
