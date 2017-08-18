package main

import (
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
		//conf.Log.Printf("File %s does not exist: %v\n", fileName, err)
		return []byte{}, err
	}

	data := make([]byte, fileStat.Size())
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		//conf.Log.Printf("Error while opening %s: %v\n", fileName, err)
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

	go gameMaker()

	splash, _ := getDataFromFile("splash.txt")
	port := flag.Int("p", 4242, "Port to listen")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("error in net.Listen : %s", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error in ln.Accept : %s", err)
		}
		fmt.Printf("%v : new gamer\n", conn)
		defer conn.Close()
		go handleConnection(conn, splash)
	}
}

func handleConnection(conn net.Conn, splash []byte) {
	player, err := getPlayerData(conn, splash)
	if err != nil {
		//
	}
	Mutex.Lock()
	Players = append(Players, &player)
	Mutex.Unlock()
}

func gameMaker() {
	for {
		if len(Players) > 1 {
			p1 := Players[0]
			p2 := Players[1]
			Mutex.Lock()
			Players = Players[2:]
			Mutex.Unlock()
			p1.inGame = true
			p2.inGame = true
			go Round(p1, p2)
		}
	}
}
