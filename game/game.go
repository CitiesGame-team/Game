package game

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"time"

	"../config"
	"../helpers"
)

const maxNameLength = 25
const maxDelay = 1000

var Mutex *sync.Mutex = &sync.Mutex{}

//To format the users output
// http://www.isthe.com/chongo/tech/comp/ansi_escapes.html
var (
	home       = []byte{27, 91, 72}
	clear      = []byte{27, 91, 50, 74}
	down       = []byte{27, 91, 1, 66}
	up         = []byte{27, 91, 65}
	colorRed   = []byte("\x1b[33m")
	colorGreen = []byte("\x1b[32m")
	colorWhite = []byte("\x1b[37m")
)

// Get data of player and return the structure
func getPlayerData(conn net.Conn, splash []byte) (*Player, error) {
	_, err := conn.Write(clear)
	if err != nil {
		return nil, errors.New("Communication error")
	}
	_, err = conn.Write(home)
	if err != nil {
		return nil, errors.New("Communication error")
	}
	_, err = conn.Write(splash)
	if err != nil {
		return nil, errors.New("Communication error")
	}

	io := bufio.NewReader(conn)

	line, err := io.ReadString('\n')
	if err != nil {
		return nil, errors.New("Communication error")
	}
	_, err = conn.Write(down)
	if err != nil {

	}
	name := strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1)
	if name == "" {
		return nil, errors.New("Empty name")
	}
	if len(name) > maxNameLength {
		return nil, errors.New("Too long name")
	}

	log.Printf("New user connected: %s\n", name)
	off := make(chan bool)
	game := make(chan string)
	return &Player{Conn: conn, Name: name, offline: off, gameChanges: game}, nil
}

func RunGame(conf config.ProjectConfig) {
	log.Printf("Starting %s...", conf.Name)

	splash, _ := helpers.ReadFile(conf.SplashFile)

	go gameMaker()

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))

	if err != nil {
		log.Fatalf("error in net.Listen : %s", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error in ln.Accept : %s", err)
		}
		go handleConnection(conn, splash)
	}
}

func handleConnection(conn net.Conn, splash []byte) {
	defer conn.Close()
	player, err := getPlayerData(conn, splash)
	if err != nil {
		log.Printf("Couldn't log in: %s\n", err.Error())
	}
	addPlayer(player)
	go player.reader()
	for {
		time.Sleep(10 * time.Millisecond)
		if player.game == nil {
			select {
			case <-player.offline:
				log.Printf("User %s disconnected.\n", player.Name)
				removePlayer(player)
			case massege := <-player.gameChanges:
				player.Conn.Write([]byte(massege))
			case <-time.After(time.Second * 120):
				log.Printf("User %s starts game with bot", player.Name)
				//go bot(player)
			}
		} else {
			select {
			case <-player.offline:
				log.Printf("User %s disconnected.\n", player.Name)
				//and something more
				return
			case command := <-player.game.chIn:
				if command == "exit" {
					player.Conn.Write(colorRed)
					player.Conn.Write([]byte(fmt.Sprintf("Your oppnent %s disconnected. You are winner.\nWait for a new opponent.\n", player.game.opponentName)))
					player.Conn.Write(colorWhite)
					player.game = nil
					Players[player] = true
				} else {
					player.Conn.Write(colorRed)
					player.Conn.Write([]byte(player.game.opponentName + ": "))
					player.Conn.Write(colorWhite)
					player.Conn.Write([]byte(command + "\n"))
					player.game.lastTown = command
				}
			case <-time.After(time.Second * 120):
				if player.game.priority != *player.game.stage {
					player.Conn.Write([]byte("You are winner.\n"))
				} else {
					player.Conn.Write([]byte("Time out. You are loser.\n"))
				}
				player.game = nil
				addPlayer(player)
			}
		}
	}
}

func addPlayer(p *Player) {
	Mutex.Lock()
	defer Mutex.Unlock()
	Players[p] = true
}

func removePlayer(p *Player) {
	Mutex.Lock()
	defer Mutex.Unlock()
	delete(Players, p)

}

func gameMaker() {
	for {
		safetyGameMaker()
		time.Sleep(10 * time.Millisecond)
	}
}

func safetyGameMaker() {
	Mutex.Lock()
	defer Mutex.Unlock()
	if len(Players) > 1 {
		i := 0
		var p1, p2 *Player
		for p := range Players {
			if i == 0 {
				p1 = p
				delete(Players, p)
			} else if i == 1 {
				p2 = p
				delete(Players, p)
			} else {
				break
			}
			i++
		}
		ch1 := make(chan string)
		ch2 := make(chan string)

		massege := string(colorGreen) +
			fmt.Sprintf("Your oponent is %s. You starts.\n", p2.Name) + string(colorWhite)
		p1.gameChanges <- massege

		massege = string(colorGreen) +
			fmt.Sprintf("Your oponent is %s. %s starts.\n", p1.Name, p1.Name) + string(colorWhite)
		p2.gameChanges <- massege

		p1.game = &Game{chIn: ch1, chOut: ch2, opponentName: p2.Name, priority: 0, stage: new(int), lastTown: ""}
		p2.game = &Game{chIn: ch2, chOut: ch1, opponentName: p1.Name, priority: 1, stage: p1.game.stage, lastTown: ""}
		log.Printf("New game: %s - %s\n", p1.Name, p2.Name)
	}
}
