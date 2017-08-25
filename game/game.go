package game

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"Game/config"
	"Game/db"
	"Game/helpers"
)

var (
	back  = []byte{27, 91, 1, 68}
	red   = []byte("\x1b[33m")
	green = []byte("\x1b[32m")
	white = []byte("\x1b[37m")
)

const maxNameLength = 25
const timeForMove = time.Minute
const timeToWaitOpponent = time.Second * 10

var Mutex *sync.Mutex = &sync.Mutex{}

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
			//it will be closed by defer
			//conn.Close()
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
				helpers.SendText(conn, []byte("Couldn't check your auth, try again!"))
				continue NAME
			}

			if pass == "" {
				// возвращаемся к вводу имени
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

		go handleConnection(conn, splash)
	}
}

func handleConnection(conn net.Conn, splash []byte) {
	defer conn.Close()
	player, err := getPlayerData(conn, splash)
	if err != nil {
		log.Printf("Couldn't log in: %s\n", err.Error())
		return
	}
	helpers.SendGreen(player.Conn, []byte("Wait for a new opponent...\n"))

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
				userModel1, userModel2 := player.userModel, player.userModel
				gameModelId, err := db.GameAdd(*userModel1, *userModel2, db.GAME_STARTED)
				if err != nil {
					return
				}
				gameModel, err := db.GameGet(gameModelId)
				if err != nil {
					return
				}

				log.Printf("User %s starts game with bot", player.Name)
				removePlayer(player)

				ch1 := make(chan string)
				ch2 := make(chan string)
				helpers.SendBack(player.Conn)
				helpers.SendGreen(player.Conn, []byte("Your opponent is bot. You start.\n"))
				helpers.SendText(player.Conn, []byte("> "))
				player.game = &Game{chIn: ch1, chOut: ch2, opponentName: "bot", priority: 0, stage: new(int), lastTown: "", gameModel: &gameModel}
				go bot(&Game{chIn: ch2, chOut: ch1, opponentName: player.Name, priority: 1, stage: player.game.stage, lastTown: "", gameModel: &gameModel})
			}
		} else {
			select {
			case <-player.offline:
				log.Printf("User %s disconnected.\n", player.Name)
				return
			case command := <-player.game.chIn:
				if command == "exit" {
					helpers.SendBack(player.Conn)
					helpers.SendBack(player.Conn)
					helpers.SendRed(player.Conn, []byte(fmt.Sprintf("Your opponent %s disconnected. You are winner.\n", player.game.opponentName)))
					helpers.SendGreen(player.Conn, []byte("Wait for a new opponent...\n"))
					//helpers.SendText(player.Conn, []byte("> "))
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
				helpers.SendBack(player.Conn)
				if player.game.priority != *player.game.stage {
					helpers.SendRed(player.Conn, []byte("You are winner.\n"))
					helpers.SendGreen(player.Conn, []byte("Wait for a new opponent...\n"))
				} else {
					helpers.SendRed(player.Conn, []byte("Time out. You are loser.\n"))
					helpers.SendGreen(player.Conn, []byte("Wait for a new opponent...\n"))
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
				if p1.userModel == nil {
					continue
				}
				delete(Players, p)
			} else if i == 1 {
				p2 = p

				if p2.userModel == nil {
					continue
				}
				delete(Players, p)
			} else {
				break
			}
			i++
		}

		userModel1, userModel2 := p1.userModel, p2.userModel
		gameModelId, err := db.GameAdd(*userModel1, *userModel2, db.GAME_STARTED)
		if err != nil {
			addPlayer(p1)
			addPlayer(p2)
			return
		}
		gameModel, err := db.GameGet(gameModelId)

		if err != nil {
			addPlayer(p1)
			addPlayer(p2)
			return
		}

		ch1 := make(chan string)
		ch2 := make(chan string)

		massege := string(green) + string(back) + string(back) +
			fmt.Sprintf("Your opponent is %s. You start.", p2.Name) + string(white) + "\n> "
		p1.gameChanges <- massege

		massege = string(green) + string(back) + string(back) +
			fmt.Sprintf("Your opponent is %s. %s starts.\n", p1.Name, p1.Name) + string(white)
		p2.gameChanges <- massege

		p1.game = &Game{chIn: ch1, chOut: ch2, opponentName: p2.Name, priority: 0, stage: new(int), lastTown: "", gameModel: &gameModel}
		p2.game = &Game{chIn: ch2, chOut: ch1, opponentName: p1.Name, priority: 1, stage: p1.game.stage, lastTown: "", gameModel: &gameModel}
	}
}
