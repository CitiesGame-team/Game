package game

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"Game/db"
	"Game/helpers"
	"unicode/utf8"
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
		message = helpers.Processing(message)

		if message == "Exit" {
			if player.game != nil {
				player.game.chOut <- "exit"
			}
			player.offline <- true
			return
		} else if message == "Stat" {
			citiesStat, err1 := helpers.GetCitiesStat()
			gameStat, err2 := helpers.GetGameStat()

			if err1 != nil && err2 != nil {
				helpers.SendRed(player.Conn, []byte("\nCannot provide statistics right now...\n\n"))
				continue
			}

			if err1 == nil {
				maxLength := 0
				for _, city := range citiesStat {
					l := utf8.RuneCountInString(city.Name)
					if maxLength < l {
						maxLength = l
					}
				}
				maxLength += 5

				helpers.SendGreen(player.Conn, []byte("Top cities:\n"))
				for _, topCity := range citiesStat {
					helpers.SendBlue(player.Conn, []byte(fmt.Sprintf("%s ", topCity.Name)))
					helpers.SendText(player.Conn, []byte(strings.Repeat(" ", maxLength-utf8.RuneCountInString(topCity.Name))))
					helpers.SendText(player.Conn, []byte(fmt.Sprintf("\t// %d\n", topCity.Count)))
				}
			}

			if err2 == nil {
				helpers.SendGreen(player.Conn, []byte("Game stats:\n"))
				helpers.SendBlue(player.Conn, []byte(fmt.Sprintf("Total games: ")))
				helpers.SendText(player.Conn, []byte(fmt.Sprintf("%d\n", gameStat.Games)))
				helpers.SendBlue(player.Conn, []byte(fmt.Sprintf("Players: ")))
				helpers.SendText(player.Conn, []byte(fmt.Sprintf("%d\n", gameStat.Players)))
				helpers.SendBlue(player.Conn, []byte(fmt.Sprintf("Cities: ")))
				helpers.SendText(player.Conn, []byte(fmt.Sprintf("%d\n\n", gameStat.Cities)))
			}
		} else if message == "" {

		} else if player.game != nil && player.game.priority == *player.game.stage {
			player.safetyMove(message)
		}
	}
}

func (player *Player) safetyMove(message string) {
	player.game.lock.Lock()
	defer player.game.lock.Unlock()
	str := player.game.lastTown
	cityModel, err := db.CityGet(message)

	if err == nil && player.game.gameModel.HasCity(cityModel) {
		helpers.SendRed(player.Conn, []byte(fmt.Sprintf("This city %q is already used in this game. Think of another city!\n", message)))
		return
	}
	exist, err := helpers.CityExists(message)
	if !exist || err != nil {
		helpers.SendRed(player.Conn, []byte(fmt.Sprintf("Unknown town. Try again.\n")))
	} else if str != "" && strings.ToLower(str[len(str)-1:]) != strings.ToLower(message[:1]) {
		helpers.SendRed(player.Conn,
			[]byte(fmt.Sprintf("Think up a city starting with the letter %s.\n", strings.ToUpper(str[len(str)-1:]))))
	} else {
		cityModel, err = db.CityGet(message)

		if err != nil {
			helpers.SendRed(player.Conn, []byte("Cannot check and save your town. Try again, please!\n"))
			return
		}
		player.game.gameModel.AddCity(cityModel)
		*(player.game.stage) = (*(player.game.stage) + 1) % 2
		player.game.chOut <- message
	}
}

func (game *Game) nextMove() {
	game.lock.Lock()
	defer game.lock.Unlock()
	*game.stage = (*game.stage + 1) % 2
}
