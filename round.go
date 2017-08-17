package main

import (
	//"fmt"
	"sync"
	//"time"
)

type Round struct {
	Players   [2]Player
	Id, State int
	sync.Mutex
}

//TODO
func botGame() {

}
