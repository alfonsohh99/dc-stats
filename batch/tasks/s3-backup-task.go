package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson"
)

func S3BackupTask(ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	var guildList []model.Guild

	guildList = []model.Guild{}

	cursor, err := database.DataCollection.Find(ctx, bson.D{})

	if err != nil {
		log.Println("ERROR FETCHING DATA GUILDS", err)
	}

	for next := cursor.Next(ctx); next; next = cursor.Next(ctx) {
		guildItem := model.Guild{}
		cursor.Decode(&guildItem)
		guildList = append(guildList, guildItem)
	}

	guildListString, err := json.Marshal(guildList)

	if err != nil {
		log.Println("Error JSON parsing guildList", err)
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Println("Error creating AWS config", err)
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	dateNow := time.Now()
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String("dc-stats"),
		Key:    aws.String("backups/" + dateNow.Format("2006-01") + "/" + dateNow.Format("2006-01-06 15:04:05")),
		Body:   strings.NewReader(string(guildListString)),
	})

	if err != nil {
		log.Println("ERROR UPLOADING BACKUP")
	} else {
		log.Println("BACKUP UPLOADED CORRECTLY")
	}

}
