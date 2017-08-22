// cities project main.go
package main

import (
	"../databases"
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Player struct {
	Conn   net.Conn
	Name   string
	Login  string
	id     int
	inGame int
	ch     chan string
}

var Players []*Player

func (player *Player) sendWait() {
	for _, r := range `|/-\` {
		player.Conn.Write(up)
		player.Conn.Write(up)
		player.Conn.Write([]byte(fmt.Sprintf("\nWaiting for opponent %c\n", r)))
		time.Sleep(100 * time.Millisecond)
	}
}

func (player *Player) getTown(str string) (string, error) {
	var town string
	for {
		player.Conn.Write(colorRed)
		player.Conn.Write([]byte(fmt.Sprintf("%s:", player.Name)))
		player.Conn.Write(colorWhite)

		io := bufio.NewReader(player.Conn)

		line, err := io.ReadString('\n')
		if err != nil {
			return "", errors.New("Communication error")
		}
		town = strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1)
		town = strings.ToUpper(town[:1]) + strings.ToLower(town[1:])

		exist, _ := databases.CheckCityDB(town)
		if !exist {
			player.Conn.Write(colorRed)
			player.Conn.Write([]byte(fmt.Sprintf("Unknown town. Try again.\n")))
		} else if str != "" && str[len(str)-2:len(str)-1] != strings.ToLower(town[:1]) {
			player.Conn.Write(colorRed)
			player.Conn.Write([]byte(fmt.Sprintf("Think up a city starts with the letter %s.\n", strings.ToUpper(str[len(str)-2:len(str)-1]))))
		} else {
			break
		}
	}
	return town, nil
}
