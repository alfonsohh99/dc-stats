package batch

import (
	"context"
	"dc-stats/batch/tasks"
	"dc-stats/constants"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/procyon-projects/chrono"
)

var (
	fetchVoiceDataTask   *chrono.ScheduledTask
	processVoiceDataTask *chrono.ScheduledTask
)

func Start(goBot *discordgo.Session) {

	taskScheduler := chrono.NewDefaultTaskScheduler()

	/**
	 * 	GATHER VOICE STATS TASK
	 */
	task, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherVoiceStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchDataInterval)
	fetchVoiceDataTask = &task

	if err == nil {
		log.Print("FetchDataTask has been scheduled successfully.  Fixed delay: ", constants.FetchDataInterval)
	}

	/**
	 * 	PROCESS VOICE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.ProcessVoiceStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.ProcessDataInterval)
	processVoiceDataTask = &task

	if err == nil {
		log.Print("processDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessDataInterval)
	}

}
