package game

import (
	"Game/db"
	"fmt"
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
			} else if len(command) > 0 {
				lastLetter := command[len(command)-1:]
				cityModel, err := db.CityHintForGame(lastLetter, *game.gameModel)

				if err != nil {
					game.chOut <- fmt.Sprintf("I do not know cities, starting with %q", lastLetter)
					game.chOut <- "exit"
					return
				}

				game.gameModel.AddCity(cityModel)
				game.lastTown = cityModel.Name
				game.chOut <- cityModel.Name
				game.nextMove()
			}
		case <-time.After(timeForMove):
			return
		}
	}
}
