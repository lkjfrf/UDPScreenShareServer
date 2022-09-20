package content

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"sync"
)

type GlobalSession struct {
	GSession sync.Map
	ConMap   sync.Map
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
	Conn    *net.UDPConn
	Addr    *net.UDPAddr
	Channel int32
}

func (gs *GlobalSession) Init() {
	gs.GSession = sync.Map{}
	gs.ConMap = sync.Map{}
}

func (gs *GlobalSession) NewPlayer(id string, c *net.UDPConn, addr *net.UDPAddr, channelNum int32) {
	Con_player := &Player{}
	Con_player.Conn = c
	Con_player.Channel = channelNum
	Con_player.Addr = addr
	gs.GSession.Store(id, Con_player)
	gs.ConMap.Store(c, id)
}

func (gs *GlobalSession) GetPlayerIdByCon(c *net.UDPConn) string {
	if id, ok := gs.ConMap.Load(c); ok {
		return id.(string)
	}
	return ""
}

func (gs *GlobalSession) BroadCastToSameChannelNumExpetMe(ChannelNum int32, Id string, recvpkt any, pkttype uint16) {
	gs.GSession.Range(func(key, value any) bool {
		if value.(*Player).Channel == ChannelNum && key != Id {
			sendBuffer := MakeSendBuffer(pkttype, recvpkt)
			gs.SendByte(value.(*Player).Conn, value.(*Player).Addr, sendBuffer)
		}
		return true
	})
}

func (gs *GlobalSession) BroadCastToSameChannelExpetMe(id string, recvpkt any, pkttype uint16) {
	var TargetChannel int32
	if p, ok := gs.GSession.Load(id); ok {
		TargetChannel = p.(*Player).Channel
	}
	gs.GSession.Range(func(key, value any) bool {
		if value.(*Player).Channel == TargetChannel && key != id {
			sendBuffer := MakeSendBuffer(pkttype, recvpkt)
			gs.SendByte(value.(*Player).Conn, value.(*Player).Addr, sendBuffer)
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

func (gs *GlobalSession) SendPacketByConn(conn *net.UDPConn, addr *net.UDPAddr, recvpkt any, pkttype uint16) {
	sendBuffer := MakeSendBuffer(pkttype, recvpkt)

	gs.SendByte(conn, addr, sendBuffer)
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
