package bts

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
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
==================================================
  ██████╗  ██████╗       ██████╗ ████████╗███████╗
 ██╔════╝ ██╔═══██╗      ██╔══██╗╚══██╔══╝██╔════╝
 ██║  ███╗██║   ██║█████╗██████╔╝   ██║   ███████╗
 ██║   ██║██║   ██║╚════╝██╔══██╗   ██║   ╚════██║
 ╚██████╔╝╚██████╔╝      ██████╔╝   ██║   ███████║
  ╚═════╝  ╚═════╝       ╚═════╝    ╚═╝   ╚══════╝
================================= lazyguyid v0.0.1
`

type (
	Tower struct {
		Name         string
		Transmitters map[string]*Transmitter
		Receiver     *Receiver
		Status       string
		ID           string
		PacketSize   int
		wg           *sync.WaitGroup
		Address      string
		Services     interface{}
	}

	// Setup bts tower
	Setup struct {
		Name       string
		TowerID    string
		TowerAddr  string
		PacketSize int
		TowerGate  bool
	}

	Receiver struct {
		Addr      *net.UDPAddr
		Conn      *net.UDPConn
		TransAddr map[string]net.UDPAddr
		Reader    func(t *Tower, v []byte) error
		Active    bool
	}

	Transmitter struct {
		TowerID      string
		TowerAddr    string
		PingInterval int64
		Active       bool
		Receiver     func(t *Tower, v []byte) error
		Conn         net.Conn
	}
)

func NewTower(setup *Setup) *Tower {
	tower := new(Tower)

	tower.Name = setup.Name
	tower.Status = "UP"
	tower.ID = setup.TowerID
	tower.wg = new(sync.WaitGroup)
	tower.Address = setup.TowerAddr
	tower.Receiver = new(Receiver)

	// setup logger

	return tower
}

func (tower *Tower) Transmitter(ts []Transmitter) {
	tower.Transmitters = make(map[string]*Transmitter, 0)
	for _, t := range ts {
		tower.Transmitters[t.TowerID] = &t
	}
}

func (tower *Tower) SetStatus(s string, towers []string) error {

	return nil
}

func (tower *Tower) Disconnected(t string) bool {
	return false
}

func (tower *Tower) Ready() error {

	fmt.Println(logo)

	// connect with transmitter
	tower.connectTransmitter()

	// create receiver
	tower.createTransmitter()

	tower.wg.Wait()

	return nil
}

func (tower *Tower) connectTransmitter() {

	for _, t := range tower.Transmitters {
		tower.wg.Add(1)
		go func() {

			defer tower.wg.Done()

			packet := make([]byte, 2048)
			protocol, host, port := convAddr(t.TowerAddr)
			conn, err := net.Dial(protocol, fmt.Sprintf("%s:%d", host, port))
			tower.Transmitters[t.TowerID].Conn = conn
			if err != nil {
				fmt.Printf("cannot making dial to tower address %s: %v", fmt.Sprintf("%s:%d", host, port), err)
			}

			defer tower.Transmitters[t.TowerID].Conn.Close()

			for tower.Transmitters[t.TowerID].Active {
				_, err = bufio.NewReader(conn).Read(packet)
				if err == nil {
					t.Receiver(tower, packet)
				} else {
					fmt.Printf("error when read the packet : %v\n", err)
				}
			}
		}()
	}
}

func (tower *Tower) createTransmitter() {
	protocol, host, port := convAddr(tower.Address)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}

	conn, err := net.ListenUDP(protocol, &addr)
	if err != nil {
		panic(err)
	}

	tower.Receiver.Conn = conn
	if err != nil {
		panic(err.Error())
	}

	log.Info("Transmitter & Receiver is Running...")
	tower.Receiver.Active = true

	tower.wg.Add(1)
	go func() {

		defer tower.wg.Done()

		for tower.Receiver.Active {
			packet := make([]byte, 2048)
			_, _, err := tower.Receiver.Conn.ReadFromUDP(packet)
			if err != nil {
				fmt.Printf("Got an error when try to read the message: %v", err)
			} else {
				if tower.Receiver.Reader != nil {
					err = tower.Receiver.Reader(tower, packet)
					if err != nil {
						fmt.Printf("Got an error when try to process the message: %v", err)
					}
				} else {
					log.Info("no reader to process the packet")
				}
				log.Info(string(packet))
			}
		}
	}()

}

func (tower *Tower) Transmit(t string, p interface{}) error {
	// _, err := fmt.Fprintf(tower.Transmitters[t].Conn, p)
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
