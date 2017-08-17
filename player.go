// cities project main.go
package main

import (
	"net"
	//"time"
)

type Player struct {
	Conn net.Conn
	Name string
}

//TODO
func (p *Player) initPlayer(id int) {
}

func (player *Player) writeToPlayer(message []byte, clean bool) {
	if clean {
		_, err := player.Conn.Write(clear)
		if err != nil {
			//kick user
		}
	}
	_, err := player.Conn.Write(message)
	if err != nil {
		//
	}
}
