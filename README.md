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

func main() {
    // create a new tower
    tower := bts.NewTower(&bts.Setup{
        Name: "My BTS",
        TowerAddr: "udp://localhost:4000",
        ID: "your-bts-ID",
        IsGateTower: true,
    })

    // making peer connection between tower
    tower.MakingPeerConnection(map[string]bts.CableLine{}{
        "near-tower-1": bts.CableLine{
            ID: "your-near-tower-id",
            TowerAddr: "udp://localhost:4000",
            PingInterval: 10 * time.Second,
            Callback: func (t *bts.Tower, packet []byte) error {
                // TODO: your logic when you get the packet
                return nil
            }
        },
    })

    // activate the tower
    tower.Up()
}
```



## **What you can do ?**
```go
// ping near tower
tower.Ping("near-tower-id")

// check if near tower disconnected
if tower.Disconnected("near-tower-id") {
    // TODO: your logic if tower disconnected
}

// send a packet to near tower
tower.SendPacket("near-tower-id", []byte("hello"))
```


### <b>NOTE**</b>

    This project is under construction, so we don't encourage you to use it in production environment!

<br/>
<b>Copyright © 2021 lazyguyid.</b>
