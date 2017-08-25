package game

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"../config"
	"../db"
	"../helpers"
)

const maxNameLength = 25
const timeForMove = time.Minute
const timeToWaitOpponent = time.Second * 10

var Mutex *sync.Mutex = &sync.Mutex{}

//To format the users output
// http://www.isthe.com/chongo/tech/comp/ansi_escapes.html
var (
	home       = []byte{27, 91, 72}
	clear      = []byte{27, 91, 50, 74}
	down       = []byte{27, 91, 1, 66}
	up         = []byte{27, 91, 65}
	back       = []byte{27, 91, 1, 68}
	colorRed   = []byte("\x1b[33m")
	colorGreen = []byte("\x1b[32m")
	colorWhite = []byte("\x1b[37m")
)

func sendWelcome(conn net.Conn, splash []byte) error {
	if err := helpers.SendClear(conn); err != nil {
		return fmt.Errorf("cannot send clear: %s", err)
	}
	if err := helpers.SendHome(conn); err != nil {
		return fmt.Errorf("cannot send home: %s", err)
	}
	if err := helpers.SendText(conn, splash); err != nil {
		return fmt.Errorf("cannot send splash: %s", err)
	}

	return nil
}

// Get data of player and return the structure
func getPlayerData(conn net.Conn, splash []byte) (*Player, error) {
	p := &Player{}

	if err := sendWelcome(conn, splash); err != nil {
		return p, err
	}

NAME:
	for {
		err := helpers.SendText(conn, []byte("Enter your name: "))

		if err != nil {
			continue
		}

		name, err := helpers.ReadString(conn)

		if err != nil {
			continue
		}

		if name == "" {
			helpers.SendRed(conn, []byte("\nName cannot be empty!\n"))
			continue
		}

		if name == "exit" {
			conn.Close()
			return p, fmt.Errorf("user decided to exit")
		}

		if len(name) > maxNameLength {
			helpers.SendRed(conn, []byte("\nName is too long!\n"))
			continue
		}

		userExists := db.UserExists(name)

		if userExists {
			helpers.SendGreen(conn, []byte(fmt.Sprintf("Welcome back, %s!\n", name)))
		} else {
			helpers.SendGreen(conn, []byte(fmt.Sprintf("You are going to register, %s!\n", name)))
		}

		for {
			helpers.SendText(conn, []byte("Your password: "))

			pass, err := helpers.ReadString(conn)

			log.Printf("User pass is %q, %s", pass, err)

			if err != nil {
				continue
			}

			if pass == "" {
				continue NAME
			}

			if userExists {
				if err = db.UserAuth(name, []byte(pass)); err != nil {
					helpers.SendRed(conn, []byte(fmt.Sprintf("Wrong password, %s!\n", name)))
					continue
				}
				break
			} else {
				created, _, err := db.UserAdd(name, []byte(pass))
				if err != nil || !created {
					continue
				}
				break
			}
		}

		userModel, err := db.UserGet(name)

		if err != nil {
			continue
		}

		log.Printf("User connected: %s\n", name)

		off := make(chan bool)
		game := make(chan string)

		return &Player{Conn: conn, Name: name, offline: off, gameChanges: game, userModel: &userModel}, nil
	}
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

		// all, err := ioutil.ReadAll(conn)
		// log.Printf("%s", all)

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
				return
			case massege := <-player.gameChanges:
				helpers.SendBack(player.Conn)
				helpers.SendText(player.Conn, []byte(massege))
			case <-time.After(timeToWaitOpponent):
				log.Printf("User %s starts game with bot", player.Name)
				ch1 := make(chan string)
				ch2 := make(chan string)
				massege := string(colorGreen) + string(back) +
					"Your opponent is bot. You start." + string(colorWhite) + "\n> "
				helpers.SendText(player.Conn, []byte(massege))
				player.game = &Game{chIn: ch1, chOut: ch2, opponentName: "bot", priority: 0, stage: new(int), lastTown: ""}
				go bot(&Game{chIn: ch2, chOut: ch1, opponentName: player.Name, priority: 1, stage: player.game.stage, lastTown: ""})
			}
		} else {
			select {
			case <-player.offline:
				log.Printf("User %s disconnected.\n", player.Name)
				return
			case command := <-player.game.chIn:
				if command == "exit" {
					helpers.SendRed(player.Conn, []byte(fmt.Sprintf("Your opponent %s disconnected. You are winner.\nWaiting for a new opponent...\n", player.game.opponentName)))
					helpers.SendText(player.Conn, []byte("> "))
					player.game = nil
					Players[player] = true
				} else {
					helpers.SendBack(player.Conn)
					helpers.SendRed(player.Conn, []byte(player.game.opponentName+": "))
					helpers.SendText(player.Conn, []byte(command+"\n> "))
					player.game.lastTown = command
				}
			case <-time.After(timeForMove):
				helpers.SendBack(player.Conn)
				if player.game.priority != *player.game.stage {
					helpers.SendRed(player.Conn, []byte("You are winner.\n> "))
				} else {
					helpers.SendRed(player.Conn, []byte("Time out. You are loser.\n> "))
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
			fmt.Sprintf("Your opponent is %s. You start.", p2.Name) + string(colorWhite) + "\n> "
		p1.gameChanges <- massege

		massege = string(colorGreen) +
			fmt.Sprintf("Your opponent is %s. %s starts.\n", p1.Name, p1.Name) + string(colorWhite)
		p2.gameChanges <- massege

		p1.game = &Game{chIn: ch1, chOut: ch2, opponentName: p2.Name, priority: 0, stage: new(int), lastTown: ""}
		p2.game = &Game{chIn: ch2, chOut: ch1, opponentName: p1.Name, priority: 1, stage: p1.game.stage, lastTown: ""}
	}
}
