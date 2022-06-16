package tasks

import (
	"context"
	botConfig "dc-stats/config"
	"dc-stats/database"
	"dc-stats/model/data"
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

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(botConfig.AWSRegion),
		config.WithCredentialsProvider(
			aws.CredentialsProviderFunc(
				func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{AccessKeyID: botConfig.AccessKeyID, SecretAccessKey: botConfig.SecretAccessKey}, nil
				})))

	if err != nil {
		log.Println("Error configurating AWS client: ", err)
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	dateNow := time.Now()

	for next := cursor.Next(ctx); next; next = cursor.Next(ctx) {
		guildItem := dataModel.Guild{}
		cursor.Decode(&guildItem)

		guildString, err := json.Marshal(guildItem)

		if err != nil {
			log.Println("Error JSON parsing guild: ", err)
			return
		}

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(botConfig.S3Bucket),
			Key:    aws.String("backups/" + guildItem.GuildID + "/" + dateNow.Format("2006-01") + "/" + dateNow.Format(time.RFC3339)),
			Body:   strings.NewReader(string(guildString)),
		})

		if err != nil {
			log.Println("ERROR UPLOADING GUILD BACKUP, ", guildItem.GuildID, err)
		}

	}

	log.Println("Finished backing up all guilds")

}
