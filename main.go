package main

import (
	"fmt"
	"vc-stats/bot"    //we will create this later
	"vc-stats/config" //we will create this later
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
