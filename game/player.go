// cities project main.go
package main

import (
	"Game/databases"
	"bufio"
	//	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Player struct {
	Conn        net.Conn
	Name        string
	Login       string
	time        int
	offline     chan bool
	gameChanges chan string
	game        *Game
}

type Game struct {
	chIn         chan string
	chOut        chan string
	opponentName string
	priority     int
	lock         sync.Mutex
	stage        *int
	lastTown     string
}

var Players map[*Player]bool = make(map[*Player]bool)

func (player *Player) sendWait() {
	for _, r := range `|/-\` {
		player.Conn.Write(up)
		player.Conn.Write(up)
		player.Conn.Write([]byte(fmt.Sprintf("\nWaiting for opponent %c\n", r)))
		time.Sleep(100 * time.Millisecond)
	}
}

func (player *Player) reader() {
	io := bufio.NewReader(player.Conn)
	for {
		time.Sleep(100 * time.Millisecond)
		message, err := io.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			player.offline <- true
			player.game.chOut <- "exit"
			return
		}
		message = strings.Replace(strings.Replace(message, "\n", "", -1), "\r", "", -1)

		if message == "exit" {
			player.offline <- true
			player.game.chOut <- "exit"
			return
		} else if message == "" {

		} else if player.game != nil && player.game.priority == *player.game.stage {

			str := player.game.lastTown
			words := strings.Split(message, " ")
			town := strings.ToUpper(words[0][:1]) + strings.ToLower(words[0][1:])
			for _, word := range words[1:] {
				town = town + " " + strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
			exist, _ := databases.CheckCityDB(town)
			if !exist {
				player.Conn.Write(colorRed)
				player.Conn.Write([]byte(fmt.Sprintf("Unknown town. Try again.\n")))
				player.Conn.Write(colorWhite)
			} else if str != "" && str[len(str)-1:] != strings.ToLower(town[:1]) {
				player.Conn.Write(colorRed)
				player.Conn.Write([]byte(fmt.Sprintf("Think up a city starts with the letter %s.\n", strings.ToUpper(str[len(str)-1:]))))
				player.Conn.Write(colorWhite)
			} else {
				player.inc()
				player.game.chOut <- town
			}
			player.Conn.Write(colorWhite)
		}
	}
}

func (player *Player) inc() {
	player.game.lock.Lock()
	defer player.game.lock.Unlock()
	*player.game.stage = (*player.game.stage + 1) % 2
}
