package game

import (
	"log"
	"time"
)

func bot(game *Game) {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case command := <-game.chIn:
			if command == "exit" {
				log.Printf("%s's bot closed.\n", game.opponentName)
				return
			} else {
				town := "Yaroslavl"
				game.lastTown = town
				game.chOut <- town
				game.nextMove()
			}
		case <-time.After(timeForMove):
			//
		}
	}
}
