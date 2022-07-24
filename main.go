package main

import (
	"context"
	"dc-stats/api"
	"dc-stats/batch"
	"dc-stats/bot"
	"dc-stats/config"
	"dc-stats/database"
	"dc-stats/logging"
	"fmt"
)

var ctx = context.TODO()

func main() {

	logging.Start()

	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	database.Start(ctx)

	bot := bot.Start(ctx)

	batch.Start(bot)

	api.Start(ctx)

}
