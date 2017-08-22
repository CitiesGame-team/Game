package main

import (
	"../config"
	"../databases"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const maxNameLength = 25
const maxDelay = 1000

var Mutex = &sync.Mutex{}

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
	//conf   Config
)

func getDataFromFile(fileName string) ([]byte, error) {
	fileStat, err := os.Stat(fileName)
	if err != nil {
		log.Printf("File %s does not exist: %v\n", fileName, err)
		return []byte{}, err
	}
	data := make([]byte, fileStat.Size())
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("Error while opening %s: %v\n", fileName, err)
		os.Exit(1)
	}
	defer f.Close()
	f.Read(data)
	return data, nil
}

// Get data of player and return the structure
func getPlayerData(conn net.Conn, splash []byte) (Player, error) {
	_, err := conn.Write(clear)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(home)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(splash)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}

	io := bufio.NewReader(conn)

	line, err := io.ReadString('\n')
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(down)
	if err != nil {

	}
	name := strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1)
	if name == "" {
		return Player{}, errors.New("Empty name")
	}
	if len(name) > maxNameLength {
		return Player{}, errors.New("Too long name")
	}

	fmt.Printf("%s\n", name)
	return Player{Conn: conn, Name: name}, nil
}

func main() {
	confFile := flag.String("config", "./config.yml", "Configuration file")
	conf, err := config.ReadProjectConfig(*confFile)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting %s...", conf.Name)

	err = databases.InitCityDB(conf.Db.ConnStr, conf.Db.PoolMaxIdle, conf.Db.PoolMaxOpen)

	go gameMaker()

	splash, _ := getDataFromFile(conf.SplashFile)

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
		//
	}
	Mutex.Lock()
	Players = append(Players, &player)
	Mutex.Unlock()

	opponentsTown := ""
	for {
		if player.inGame == 1 {
			defer close(player.ch)
			town, _ := player.getTown(opponentsTown)
			player.ch <- "\x1b[32m" + player.Name + "\x1b[37m: " + town + "\n"
			opponentsTown, _ = <-player.ch
			/*if something_bad {
				opponentsTown = ""
				player.inGame = 0
				close(player.ch)
				player.ch = nil
				Mutex.Lock()
				Players = append(Players, &player)
				Mutex.Unlock()
				fmt.Println("communication problems\n")
			} else {
			*/
			player.Conn.Write([]byte(opponentsTown))
			//}
		} else if player.inGame == 2 {
			defer close(player.ch)
			opponentsTown, _ := <-player.ch
			fmt.Println(opponentsTown[len(opponentsTown):])
			/*if something_bad {
				opponentsTown == ""
				player.inGame = 0
				close(player.ch)
				player.ch = nil
				Mutex.Lock()
				Players = append(Players, &player)
				Mutex.Unlock()
				fmt.Println("communication problems\n")
			} else {
			*/
			player.Conn.Write([]byte(opponentsTown))
			town, _ := player.getTown(opponentsTown)
			player.ch <- "\x1b[32m" + player.Name + "\x1b[37m: " + town + "\n"
			//}
		} else {
			player.sendWait()
		}
	}
}

func gameMaker() {
	for {
		if len(Players) > 1 {
			p1 := Players[0]
			p2 := Players[1]
			Mutex.Lock()
			Players = Players[2:]
			Mutex.Unlock()
			ch := make(chan string)
			p1.ch = ch
			p2.ch = ch
			p1.inGame = 1
			p2.inGame = 2

		}
	}
}
