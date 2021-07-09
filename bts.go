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
		Name       string
		CableLines map[string]CableLine
		Status     string
		ID         string
		PacketSize int
	}

	// Setup bts tower
	Setup struct {
		Name       string
		TowerID    string
		TowerAddr  string
		PacketSize int
		TowerGate  bool
	}

	CableLine struct {
		TowerID      string
		TowerAddr    string
		PingInterval int64
		Active       bool
		Callback     func(t *Tower, v []byte) error
	}
)

func NewTower(setup *Setup) *Tower {
	tower := new(Tower)

	tower.Name = setup.Name
	tower.Status = "UP"
	tower.ID = setup.TowerID

	return tower
}

func (tower *Tower) MakingPeerConnection(cablelines map[string]CableLine) {
	var wg sync.WaitGroup
	tower.CableLines = cablelines
	for _, cableline := range cablelines {
		wg.Add(1)
		go tower.connect(&cableline, &wg)
	}

	wg.Wait()
}

func (tower *Tower) SetStatus(s string, towers []string) error {
	return nil
}

func (tower *Tower) SendPacket(i interface{}, towers []string) error {
	return nil
}

func (tower *Tower) Disconnected(t string) bool {
	return false
}

func (tower *Tower) Activate() error {
	return nil
}

func (tower *Tower) connect(c *CableLine, wg *sync.WaitGroup) {
	packet := make([]byte, 2048)
	protocol, host, port := convAddr(c.TowerAddr)
	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("cannot making dial to tower address %s: %v", fmt.Sprintf("%s:%d", host, port), err)
	}

	defer wg.Done()
	defer conn.Close()

	for tower.CableLines[c.TowerID].Active {
		_, err = bufio.NewReader(conn).Read(packet)
		if err == nil {
			c.Callback(tower, packet)
		} else {
			fmt.Printf("error when read the packet : %v\n", err)
		}
	}
}

func (tower *Tower) SendSignal()

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
