package batch

import (
	"context"
	"log"
	"sync"
	"vc-stats/batch/tasks"
	"vc-stats/constants"

	"github.com/bwmarrin/discordgo"
	"github.com/procyon-projects/chrono"
)

var (
	fetchDataTask   *chrono.ScheduledTask
	processDataTask *chrono.ScheduledTask
)

func Start(goBot *discordgo.Session) {

	taskScheduler := chrono.NewDefaultTaskScheduler()

	task, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchDataInterval)
	fetchDataTask = &task

	if err == nil {
		log.Print("FetchDataTask has been scheduled successfully.  Fixed delay: ", constants.FetchDataInterval)
	}

	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.ProcessStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.ProcessDataInterval)
	processDataTask = &task

	if err == nil {
		log.Print("processDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessDataInterval)
	}

}
