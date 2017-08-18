// cities project main.go
package main

import (
	"fmt"
	"net"
	"time"
)

type Player struct {
	Conn   net.Conn
	Name   string
	Login  string
	id     int
	inGame bool
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

func (player *Player) sendWait() {
	for _, r := range `|/-\` {
		player.Conn.Write(up)
		player.Conn.Write([]byte(fmt.Sprintf("\nWaiting for opponent %c\n", r)))
		time.Sleep(100 * time.Millisecond)
		player.Conn.Write(up)
	}
}
