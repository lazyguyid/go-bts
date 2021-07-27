# 🗼 **go-bts** 
<p align="center"><img src="https://raw.githubusercontent.com/lazyguyid/go-bts/main/logo.png" align="center" /></p>
<p align="center">🗼 ~~ 🗼 ~~ 🗼</p>
<p align="center">~&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~~~&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~ </p>
<p align="center">🗼 ~~~~~~ 🗼 ~~~~~~🗼</p>
<p align="center">~&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~~~&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~ </p>

<p align="center">🗼 ~~ 🗼 ~~ 🗼</p>
<br>
it's just a simple peer to peer library to help you making a simple connection between your application/service just like a BTS Tower base on UDP network.


<br/>

## **Example**

```go
package main


import (
    bts "github.com/lazyguyid/go-bts"
)

var receiver, transmitter bool

func init() {
	flag.BoolVar(&receiver, "receiver", false, "activate receiver")
	flag.BoolVar(&transmitter, "transmitter", false, "activate transmitter")
}

func receivers(t *bts.Tower) {
	fmt.Println("receiver start")
	t = bts.NewTower(&bts.Setup{
		Name:       "Tower A",
		ID:         "bts-example",
		Addr:       "udp://127.0.0.1:4321",
		PacketSize: 2048,
		AsGate:     false,
	})
}

func transmitters(t *bts.Tower) {
	t.Connect([]bts.Transmitter{
		bts.Transmitter{
			Active:       true,
			ID:           "uuid",
			Addr:         "udp://127.0.0.1:4321",
			PingInterval: 10 * time.Second,
			Receiver: func(t *bts.Tower, v []byte, transmitter *bts.Transmitter) error {
				fmt.Println(fmt.Sprintf("\r[%s]:: %s", transmitter.Conn.RemoteAddr().String(), string(v)))
				return nil
			},
		},
	})
}

func main() {
	flag.Parse()
	tower := bts.NewTower(nil)
	if receiver {
		tower = bts.NewTower(&bts.Setup{
			Name:       "Tower A",
			ID:         "bts-example",
			Addr:       "udp://127.0.0.1:4321",
			PacketSize: 2048,
			AsGate:     false,
			Callback: func(t *bts.Tower, p []byte, transAddr net.Addr) error {
				fmt.Println(fmt.Sprintf("\r[%v]:: %s", transAddr, string(p)))
				return nil
			},
		})
	}
	if transmitter {
		transmitters(tower)
	}
	tower.ActivatePrompt(true)
	tower.Ready()
}

```
## **Screenshoot**

<p align="center"><img src="https://raw.githubusercontent.com/lazyguyid/go-bts/main/demo.gif" align="center" /></p>



### <b>NOTE**</b>

    This project is under construction, so we don't encourage you to use it in production environment!

<br/>
<b>Copyright © 2021 lazyguyid.</b>
