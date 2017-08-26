package game

import (
	"bufio"
	//	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"Game/db"
	"Game/helpers"
)

type Player struct {
	Conn        net.Conn
	Name        string
	Login       string
	time        int
	offline     chan bool
	gameChanges chan string
	game        *Game

	userModel *db.UserModel
}

type Game struct {
	chIn         chan string
	chOut        chan string
	opponentName string
	priority     int
	lock         sync.Mutex
	stage        *int
	lastTown     string

	gameModel *db.GameModel
}

var Players map[*Player]bool = make(map[*Player]bool)

func (player *Player) sendWait() {
	helpers.SendDown(player.Conn)
	for _, r := range `|/-\` {
		helpers.SendUp(player.Conn)
		player.Conn.Write([]byte(fmt.Sprintf("Waiting for opponent %c\n", r)))
		time.Sleep(100 * time.Millisecond)
	}
}

func (player *Player) reader() {
	io := bufio.NewReader(player.Conn)
	for {
		if player.game != nil && player.game.priority == *player.game.stage {
			helpers.SendText(player.Conn, []byte("> "))
		}
		message, err := io.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			if player.game != nil {
				player.game.chOut <- "exit"
			}
			player.offline <- true
			return
		}
		message = strings.Replace(strings.Replace(strings.Replace(message, "\n", "", -1), "\r", "", -1), "\b", "", -1)

		if message == "exit" {
			if player.game != nil {
				player.game.chOut <- "exit"
			}
			player.offline <- true
			return
		} else if message == "" {

		} else if player.game != nil && player.game.priority == *player.game.stage {

			str := player.game.lastTown
			words := strings.Split(message, " ")
			town := strings.ToUpper(words[0][:1]) + strings.ToLower(words[0][1:])
			for _, word := range words[1:] {
				town = town + " " + strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}

			cityModel, err := db.CityGet(town)

			if err == nil && player.game.gameModel.HasCity(cityModel) {
				helpers.SendRed(player.Conn, []byte(fmt.Sprintf("This city %q is already used in this game. Think of another city!\n", town)))
				continue
			}

			exist, err := helpers.CityExists(town)
			if !exist || err != nil {
				helpers.SendRed(player.Conn, []byte(fmt.Sprintf("Unknown town. Try again.\n")))
			} else if str != "" && strings.ToLower(str[len(str)-1:]) != strings.ToLower(town[:1]) {
				helpers.SendRed(player.Conn,
					[]byte(fmt.Sprintf("Think up a city starting with the letter %s.\n", strings.ToUpper(str[len(str)-1:]))))
			} else {
				cityModel, err = db.CityGet(town)

				if err != nil {
					helpers.SendRed(player.Conn, []byte("Cannot check and save your town. Try again, please!\n"))
					continue
				}
				player.game.gameModel.AddCity(cityModel)
				player.game.nextMove()
				player.game.chOut <- town
			}
		}
	}
}

func (game *Game) nextMove() {
	game.lock.Lock()
	defer game.lock.Unlock()
	*game.stage = (*game.stage + 1) % 2
}
