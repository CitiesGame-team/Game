package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	//"time"
)

var Players []*Player

func (player *Player) getTown() (string, error) {
	player.Conn.Write(colorRed)
	player.writeToPlayer([]byte(fmt.Sprintf("%s:", player.Name)), false)
	player.Conn.Write(colorWhite)

	io := bufio.NewReader(player.Conn)

	line, err := io.ReadString('\n')
	if err != nil {
		return "", errors.New("Communication error")
	}
	town := strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1)
	return town, nil
}

func (p *Player) sendTown(town string, name string) {
	p.Conn.Write(colorGreen)
	p.Conn.Write(([]byte(fmt.Sprintf("%s: ", name))))
	p.Conn.Write(colorWhite)
	p.Conn.Write([]byte(fmt.Sprintf("%s\n", town)))
}

func Round(p1 *Player, p2 *Player) {
	p1.writeToPlayer([]byte(fmt.Sprintf("Your opponent is %s.\nYou starts game.\n ", p2.Name)), false)
	p2.writeToPlayer([]byte(fmt.Sprintf("Your opponent is %s.\nWait.\n", p1.Name)), false)
	for {
		town1, _ := p1.getTown()
		p2.sendTown(town1, p1.Name)
		town2, _ := p2.getTown()
		p1.sendTown(town2, p2.Name)
	}
}

//TODO
func botGame() {

}
