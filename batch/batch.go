package batch

import (
	"context"
	"dc-stats/batch/tasks"
	"dc-stats/constants"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/procyon-projects/chrono"
)

var (
	fetchVoiceDataTask     *chrono.ScheduledTask
	processVoiceDataTask   *chrono.ScheduledTask
	fetchMessageDataTask   *chrono.ScheduledTask
	processMessageDataTask *chrono.ScheduledTask
	fetchNicknamesTask     *chrono.ScheduledTask
	fetchChannelNamesTask  *chrono.ScheduledTask
	processS3BackupTask    *chrono.ScheduledTask
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
	}, constants.FetchVoiceDataInterval, chrono.WithTime(time.Now().Add(constants.FetchVoiceDataInterval)))
	fetchVoiceDataTask = &task

	if err == nil {
		log.Print("FetchVoiceDataTask has been scheduled successfully.  Fixed delay: ", constants.FetchVoiceDataInterval)
	} else {
		log.Println("Error scheduling task FetchVoiceDataTask: ", err)
	}

	/**
	 * 	PROCESS VOICE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.ProcessVoiceStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.ProcessVoiceDataInterval, chrono.WithTime(time.Now().Add(constants.ProcessVoiceDataInterval)))
	processVoiceDataTask = &task

	if err == nil {
		log.Print("ProcessVoiceDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessVoiceDataInterval)
	} else {
		log.Println("Error scheduling task ProcessVoiceDataTask: ", err)
	}

	/**
	 *  GATHER MESSAGE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherMessageStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchMessageDataInterval, chrono.WithTime(time.Now().Add(constants.FetchMessageDataInterval)))
	fetchMessageDataTask = &task

	if err == nil {
		log.Print("FetchMessageDataTask has been scheduled successfully. Fixed delay: ", constants.FetchMessageDataInterval)
	} else {
		log.Println("Error scheduling task FetchMessageDataTask: ", err)
	}

	/**
	 * 	PROCESS MESSAGE STATS TASK
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.ProcessMessageStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.ProcessMessageDataInterval, chrono.WithTime(time.Now().Add(constants.ProcessMessageDataInterval)))
	processMessageDataTask = &task

	if err == nil {
		log.Print("ProcessMessageDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessMessageDataInterval)
	} else {
		log.Println("Error scheduling task ProcessMessageDataTask: ", err)
	}

	/**
	 * 	GATHER USER NICKNAMES
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherNicknameStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchNicknamesInterval, chrono.WithTime(time.Now()))
	fetchNicknamesTask = &task

	if err == nil {
		log.Print("FetchNicknamesTask has been scheduled successfully.  Fixed delay: ", constants.FetchNicknamesInterval)
	} else {
		log.Println("Error scheduling task FetchNicknamesTask: ", err)
	}

	/**
	 * 	GATHER CHANNEL NAMES
	 */
	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.GatherChannelNameStats(goBot, ctx, &wg)
		wg.Wait()
	}, constants.FetchChannelNamesInterval, chrono.WithTime(time.Now()))
	fetchChannelNamesTask = &task

	if err == nil {
		log.Print("FetchChannelNamesTask has been scheduled successfully.  Fixed delay: ", constants.FetchChannelNamesInterval)
	} else {
		log.Println("Error scheduling task FetchChannelNamesTask: ", err)
	}

	/**
	 * 	S3 BACKUP (0 0 * * *)
	 */
	task, err = taskScheduler.ScheduleWithCron(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go tasks.S3BackupTask(ctx, &wg)
		wg.Wait()
	}, constants.S3BackupCronExpression)
	processS3BackupTask = &task

	if err == nil {
		log.Print("S3BackupTask has been scheduled successfully.  Cron: ", constants.S3BackupCronExpression)
	} else {
		log.Println("Error scheduling task S3BackupTask: ", err)
	}

}
