package tasks

import (
	"context"
	botConfig "dc-stats/config"
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

	cursor, err := database.DataCollection.Find(ctx, bson.D{})

	if err != nil {
		log.Println("ERROR FETCHING DATA GUILDS", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Println("Error configurating AWS client: ", err)
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	dateNow := time.Now()

	for next := cursor.Next(ctx); next; next = cursor.Next(ctx) {
		guildItem := model.Guild{}
		cursor.Decode(&guildItem)

		guildString, err := json.Marshal(guildItem)

		if err != nil {
			log.Println("Error JSON parsing guild: ", err)
			return
		}

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(botConfig.S3Bucket),
			Key:    aws.String("backups/" + guildItem.GuildID + "/" + dateNow.Format("2006-01") + "/" + dateNow.Format("2006-01-06 15:04:05")),
			Body:   strings.NewReader(string(guildString)),
		})

		if err != nil {
			log.Println("ERROR UPLOADING BACKUP, ", guildItem.GuildID)
		}

	}

	log.Println("Finished backing up all guilds")

}
