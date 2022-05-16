package main

import (
	"context"
	"fmt"
	"vc-stats/batch"
	"vc-stats/bot"
	"vc-stats/config"
	"vc-stats/database"
)

var ctx = context.TODO()

func main() {

	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	database.Start(ctx)

	bot := bot.Start(ctx)

	batch.Start(bot)

	<-make(chan struct{})
	return
}
