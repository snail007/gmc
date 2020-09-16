package main

import (
	"github.com/snail007/gmc/timeutil"
	"log"
)

func main() {
	//Get the current time in seconds
	second := timeutil.GetNowSecond()
	log.Print(second)
}
