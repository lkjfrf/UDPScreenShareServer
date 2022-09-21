package core

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/lkjfrf/content"
)

type NetworkCore struct {
}

var instance *NetworkCore
var once sync.Once

func GetNetworkCore() *NetworkCore {
	once.Do(func() {
		instance = &NetworkCore{}
	})
	return instance
}

func (nc *NetworkCore) Init(port string) {
	log.Println("INIT_NetworkCore")
	nc.Connect(port)
}

func (nc *NetworkCore) ParseHeader(header []byte) (int, int) {
	pktsize := binary.LittleEndian.Uint16(header[:2])
	pktid := binary.LittleEndian.Uint16(header[2:4])

	return int(pktsize), int(pktid)
}

func (nc *NetworkCore) Connect(port string) {

	ServerAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("listening on ", port)

	conn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Println("Connect Fail : ", err)
	} else {
		log.Println("Connect Success : ", conn)
	}

	MaxGoRoutine := 500
	waitChan := make(chan struct{}, MaxGoRoutine)
	count := 0
	for {
		waitChan <- struct{}{}
		count++
		go func(count int) {
			nc.Recv(conn)
			<-waitChan
		}(count)
	}
	//go nc.Recv(conn)
}

func (nc *NetworkCore) Recv(conn *net.UDPConn) {
	for {
		data := make([]byte, 64*1024)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			break
		}
		//_, err = conn.WriteToUDP(data, addr)
		//log.Println("SendPacket ", addr, "/", conn, "/", string(data))

		if n > 0 && err == nil {
			pktsize, pktid := nc.ParseHeader(data)
			log.Println("RecvPacket : ", addr, " - ", "pktid : ", pktid)

			if pktsize > 4 {
				data = data[4:pktsize]
				//n, _, _ := conn.ReadFrom(recv)
				if n > 4 || pktid < 100 {
					if content.GetContentManager().HandlerFunc[(int)(pktid)] != nil {
						content.GetContentManager().HandlerFunc[(int)(pktid)](conn, addr, string(data))
					}
				}
			}
		}
	}
}
