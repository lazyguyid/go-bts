package bts

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
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

type (
	Tower struct {
		Name         string
		Transmitters map[string]*Transmitter
		Receiver     *Receiver
		Status       string
		ID           string
		PacketSize   int
		wg           sync.WaitGroup
		Address      string
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
		Addr *net.UDPAddr
		Conn *net.UDPConn
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
	tower.wg = sync.WaitGroup

	return tower
}

func (tower *Tower) BuildCableLines(cablelines []CableLine) {
	tower.CableLines = make(map[string]*CableLine, 0)
	for _, cableLine := range cablelines {
		tower.CableLines[cableLine.ID] = &cableLine
	}
}

func (tower *Tower) SetStatus(s string, towers []string) error {

	return nil
}

func (tower *Tower) Disconnected(t string) bool {
	return false
}

func (tower *Tower) StandUp() error {

	// connect with transmitter
	f.connectTransmitter()

	// create receiver
	f.createReceiver()

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
			tower.CableLine[c.TowerID].Conn, err = net.Dial(protocol, fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				fmt.Printf("cannot making dial to tower address %s: %v", fmt.Sprintf("%s:%d", host, port), err)
			}

			defer tower.CableLine[c.TowerID].Conn.Close()

			for tower.CableLines[c.TowerID].Active {
				_, err = bufio.NewReader(conn).Read(packet)
				if err == nil {
					c.Receiver(tower, packet)
				} else {
					fmt.Printf("error when read the packet : %v\n", err)
				}
			}
		}()
	}
}

func (tower *Tower) createTransmitter() {
	packet := make([]byte, 2048)
	protocol, host, port := convAddr(c.TowerAddr)
	tower.Addr = net.UDPAddr{
		Port: port,
		IP: host,
	}

	tower.Receiver.Conn, err := net.ListenUDP(protocol, tower.Addr)
	if err != nil {
		panic(err.Error())
	}

	for tower.Receiver.Active {

	}
}

func (tower *Tower) Transmit(t string, p []byte) error {
	_, err := fmt.Fprintf(tower.cableLines[t].Conn, p)
	return err
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
