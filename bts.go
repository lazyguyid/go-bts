package bts

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

/*

// create new tower
tower := bts.NewTower(&bts.Setup{
	Name: "my-bts",
	TowerAddr: "udp://localhost:4000",
	ID: "g3Lenk123",
	IsGateTower: true,
})

tower.MakingPeerConnection(map[string]bts.CableLine{}{
	"order-tower": bts.CableLine{
		ID: "my-near-bts",
		TowerAddr: "udp://localhost:4000",
		PingInterval: 10 * time.Second,
		Callback: func (t *bts.Tower, msg []byte) error {
			return nil
		}
	},
})

// set status method
tower.SetStatus(bts.Disconnected, ["baracuda"])
tower.Ping(["pyro"])

//
// handle if
if towers.Disconnected(["baracuda", "pyro"]) {
	tower.SendTowerStatus(bts.TowerDisconnect, [""])
}

tower.Up()

*/

var logo string = `

  ██████╗  ██████╗       ██████╗ ████████╗███████╗
 ██╔════╝ ██╔═══██╗      ██╔══██╗╚══██╔══╝██╔════╝
 ██║  ███╗██║   ██║█████╗██████╔╝   ██║   ███████╗
 ██║   ██║██║   ██║╚════╝██╔══██╗   ██║   ╚════██║
 ╚██████╔╝╚██████╔╝      ██████╔╝   ██║   ███████║
  ╚═════╝  ╚═════╝       ╚═════╝    ╚═╝   ╚══════╝
==================================================
= Simple, Fast, Secure ===========================
=============================== © lazyguyid %v

`
var version string = "v0.0.1"

type (
	PacketHeader struct {
		ID        string
		Timestamp int64
		Ping      bool
		Info      *HeaderInfo
	}

	Packet struct {
		Data interface{}
	}

	HeaderInfo struct{}

	Prompt struct {
		Active    bool
		CurrTrans []string
	}

	Tower struct {
		ID            string
		Name          string
		Status        string
		Transmitters  map[string]*Transmitter
		Receiver      *Receiver
		MaxBufferSize int64
		wg            *sync.WaitGroup
		Address       string
		EnablePrompt  bool
		Prompt        *Prompt
	}

	// Setup bts tower
	Setup struct {
		Name       string
		ID         string
		Addr       string
		PacketSize int
		AsGate     bool
		Callback   func(t *Tower, v []byte, transAddr net.Addr) error
	}

	Receiver struct {
		Addr      *net.Addr
		Conn      net.PacketConn
		TransAddr map[string]net.Addr
		Action    func(t *Tower, v []byte, transAddr net.Addr) error
		Active    bool
	}

	Transmitter struct {
		ID            string
		Addr          string
		PingInterval  time.Duration
		EnablePrompt  bool
		Active        bool
		Receiver      func(t *Tower, v []byte, trans *Transmitter) error
		Conn          net.Conn
		MaxBufferSize int64
	}
)

func NewTower(setup *Setup) *Tower {
	tower := new(Tower)
	tower.wg = new(sync.WaitGroup)
	if setup == nil {
		return tower
	}

	tower.Name = setup.Name
	tower.Status = "UP"
	tower.ID = setup.ID
	tower.Address = setup.Addr
	// receiver
	tower.Receiver = new(Receiver)
	tower.Receiver.Action = setup.Callback
	tower.Receiver.TransAddr = make(map[string]net.Addr)

	return tower
}

func (tower *Tower) Connect(ts []Transmitter) {
	tower.Transmitters = make(map[string]*Transmitter, 0)
	for _, t := range ts {
		tower.Transmitters[t.ID] = &t
	}
}

func (tower *Tower) SetStatus(s string, towers []string) error {

	return nil
}

func (tower *Tower) Disconnected(t string) bool {
	return false
}

func (tower *Tower) Ready() (err error) {

	// print current version
	fmt.Printf(logo, version)
	// show input prompt if it's enable
	tower.prompt()
	// connect with transmitter
	tower.connectTransmitter()
	if tower.Receiver != nil {
		// create receiver
		tower.createTransmitter()
	}
	// wait all process
	tower.wg.Wait()

	return
}

func (tower *Tower) connectTransmitter() {
	var buffSize int64 = 2048
	if tower.MaxBufferSize != 0 {
		buffSize = tower.MaxBufferSize
	}

	for _, t := range tower.Transmitters {
		tower.wg.Add(1)

		go func() {
			defer tower.wg.Done()
			protocol, host, port := convAddr(t.Addr)
			conn, err := net.Dial(protocol, fmt.Sprintf("%s:%d", host, port))
			tower.Transmitters[t.ID].Conn = conn
			if err != nil {
				fmt.Printf("can't making dial to transmitter with address %s: %v", fmt.Sprintf("%s:%d", host, port), err)
			}

			defer tower.Transmitters[t.ID].Conn.Close()
			fmt.Printf("\r=> transmitter connected with %s:%d\n\n", host, port)

			for tower.Transmitters[t.ID].Active {
				tower.printPrompt()

				packet := make([]byte, buffSize)
				_, err = bufio.NewReader(tower.Transmitters[t.ID].Conn).Read(packet)
				if err == nil {
					t.Receiver(tower, packet, tower.Transmitters[t.ID])
				} else {
					fmt.Printf("error when read the packet : %v\n", err)
				}
			}
		}()
	}
}

func (tower *Tower) ActivatePrompt(v bool) {
	tower.EnablePrompt = v
	if tower.Prompt == nil {
		tower.Prompt = new(Prompt)
		tower.Prompt.Active = true
		tower.Prompt.CurrTrans = []string{"#all"}
	}

}

func (tower *Tower) RunCmd(t string) (cntinue bool, err error) {
	arrTxt := strings.Split(t, " ")
	if len(arrTxt) <= 1 {
		return false, nil
	}

	command := arrTxt[0]
	// remove new line char
	values := strings.Replace(arrTxt[1], "\n", "", -1)

	switch command {
	case CONNECTWITH:
		if values == "all" {
			tower.Prompt.CurrTrans = []string{"#all"}
			return true, nil
		}
		var addrs []string
		for _, addr := range strings.Split(values, ",") {
			addrs = append(addrs, addr)
		}
		tower.Prompt.CurrTrans = addrs
		return true, nil
	}

	return false, nil
}

func (tower *Tower) prompt() {

	if !tower.EnablePrompt {
		return
	}

	tower.wg.Add(1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n----------------")
		fmt.Println("| Shell Active |")
		fmt.Println("----------------")

		defer tower.wg.Done()

		for {
			tower.printPrompt()
			text, _ := reader.ReadString('\n')

			// convert CRLF to LF
			isContinue, err := tower.RunCmd(text)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			// pass the logic if isContinue == true
			if isContinue {
				continue
			}

			text = strings.Replace(text, "\n", "", -1)
			currAddr := tower.Prompt.CurrTrans
			if len(currAddr) > 0 && currAddr[0] == "#all" {
				// send it to transmitter & receiver
				go func() {
					for _, transmitter := range tower.Transmitters {
						go fmt.Fprintf(transmitter.Conn, text)
					}
				}()
				if tower.Receiver != nil && tower.Receiver.Active {
					go func() {
						for _, addr := range tower.Receiver.TransAddr {
							_, err := tower.Receiver.Conn.WriteTo([]byte(text), addr)
							if err != nil {
								log.Warn(err.Error())
							}
						}
					}()
				}
			} else {
				for _, addrID := range currAddr {
					tower.Receiver.Conn.WriteTo(
						[]byte(text),
						tower.Receiver.TransAddr[addrID],
					)
				}
			}
		}
	}()
}

func (tower *Tower) printPrompt() {
	if !tower.EnablePrompt {
		return
	}

	if tower.Prompt != nil && tower.Prompt.Active {
		actvTrans := tower.Prompt.CurrTrans
		if len(actvTrans) > 0 && actvTrans[0] != "#all" {
			fmt.Print(fmt.Sprintf("[%s]:: ", strings.Join(actvTrans, ",")))
		} else {
			fmt.Print("[#all]:: ")
		}

	} else {
		fmt.Print("[#all]:: ")
	}

}

func (tower *Tower) createTransmitter() {
	protocol, host, port := convAddr(tower.Address)
	conn, err := net.ListenPacket(protocol, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	tower.Receiver.Conn = conn
	tower.Receiver.Active = true
	log.Info("\r=> receiver is running...\n\n")
	tower.wg.Add(1)
	go func() {
		defer tower.wg.Done()
		for tower.Receiver.Active {
			// show input prompt if prompt is enable
			tower.printPrompt()
			packet := make([]byte, 2048)
			_, transAddr, err := tower.Receiver.Conn.ReadFrom(packet)
			tower.Receiver.TransAddr[transAddr.String()] = transAddr
			// if !tower.Receiver.IsValidTransmitter(transAddr) {
			// 	log.Warn("invalid message from anonymous detected")
			// }
			if err != nil {
				log.Warn("got an error when try to read the message: %v", err)
			} else {
				if tower.Receiver.Action != nil {
					err = tower.Receiver.Action(tower, packet, transAddr)
					if err != nil {
						fmt.Printf("Got an error when try to process the message: %v\n", err)
					}
				}
			}
		}
	}()

}

func (tower *Tower) Transmit(t string, p interface{}) error {
	return nil
}

func convAddr(addr string) (protocol string, host string, port int) {
	s := strings.Split(addr, ":")
	if len(s) < 3 {
		panic("error parse address")
	}

	protocol = s[0]
	host = strings.Split(s[1], "//")[1]
	port, err := strconv.Atoi(s[2])
	if err != nil {
		panic(err.Error())
	}

	return
}

// transmit to transmitter
func (transmitter *Transmitter) Transmit(v interface{}) (err error) {

	packet := new(Packet)
	packet.Data = v

	// marshal with msgpack
	bPacket, err := msgpack.Marshal(packet)
	if err != nil {
		return
	}
	// send packet to transmitter
	_, err = transmitter.Conn.Write(bPacket)
	if err != nil {
		return
	}
	return
}

// receiver
func (receiver *Receiver) IsValidTransmitter(transAddr net.Addr) (r bool) {

	return false
}
