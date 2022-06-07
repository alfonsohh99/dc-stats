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
	fetchMessageDataTask *chrono.ScheduledTask
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
	}, constants.FetchVoiceDataInterval)
	fetchVoiceDataTask = &task

	if err == nil {
		log.Print("FetchVoiceDataTask has been scheduled successfully.  Fixed delay: ", constants.FetchVoiceDataInterval)
	}

	/**
	 * 	PROCESS VOICE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.ProcessVoiceStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.ProcessVoiceDataInterval)
	processVoiceDataTask = &task

	if err == nil {
		log.Print("processVoiceDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessVoiceDataInterval)
	}

	/**
	 *  GATHER MESSAGE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherMessageStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchMessageDataInterval)
	fetchMessageDataTask = &task

	if err == nil {
		log.Print("FetchMessageDataTask has been scheduled successfully. Fixed delay: ", constants.FetchMessageDataInterval)
	}

}
