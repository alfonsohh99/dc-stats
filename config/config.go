package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token            string
	BotPrefix        string
	DatabasePassword string
	DatabaseUser     string
	DatabaseEndpoint string
	DatabasePort     string

	config *configStruct
)

type configStruct struct {
	Token            string `json : "Token"`
	BotPrefix        string `json : "BotPrefix"`
	DatabasePassword string `json : "DatabasePassword"`
	DatabaseUser     string `json : "DatabaseUser"`
	DatabaseEndpoint string `json : "DatabaseEndpoint"`
	DatabasePort     string `json : "DatabasePort"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	DatabasePassword = config.DatabasePassword
	DatabaseUser = config.DatabaseUser
	DatabaseEndpoint = config.DatabaseEndpoint
	DatabasePort = config.DatabasePort

	return nil

}
