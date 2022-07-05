package api

import (
	"context"
	"dc-stats/api/api-model"
	"dc-stats/config"
	"dc-stats/database"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx context.Context

func Start(appContext context.Context) {

	ctx = appContext

	e := echo.New()

	e.GET("auth/:code", authentication)

	e.GET("user", getUser)

	e.GET("user/guilds", getGuilds)

	e.GET("user/guilds/:guildId", getGuild)

	e.Logger.Error(e.Start(":8080"))
}

func authentication(c echo.Context) error {
	code := c.Param("code")

	authentication, err := retrieveCredentialsFromCode(code)
	if err != nil {
		log.Println("Error retrieving token: ", err)
		return c.String(http.StatusBadRequest, "Error retrieving token")
	}

	accessToken := authentication.AccessToken

	userInfo, err := getUserInfoFromAccessToken(accessToken)

	if err != nil {
		log.Println("Error retrieving user info: ", err)
		return c.String(http.StatusBadRequest, "Error retrieving user info")
	}

	user := apiModel.User{
		ID:            primitive.NewObjectID(),
		UserId:        userInfo.ID,
		UserName:      userInfo.Username,
		AccentColor:   userInfo.AccentColor,
		Discriminator: userInfo.Discriminator,
		Verified:      userInfo.Verified,
		Avatar:        userInfo.Avatar,
		PremiumType:   userInfo.PremiumType,
		Banner:        userInfo.Banner,
		Locale:        userInfo.Locale,
		MFAEnabled:    userInfo.MFAEnabled}
	database.SaveOrUpdatApiUser(user, ctx)

	response, err := json.Marshal(apiModel.UserAuth{
		UserId:       userInfo.ID,
		AccessToken:  authentication.AccessToken,
		RefreshToken: authentication.RefreshToken,
		ExpiresIn:    authentication.ExpiresIn})

	if err != nil {
		log.Println("Error creating response:  ", err)
		return c.String(http.StatusInternalServerError, "Error creating response")
	}

	return c.String(http.StatusOK, string(response))
}

func getUser(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")

	userInfo, err := getUserInfoFromAccessToken(token)
	if err != nil {
		log.Println("Error retrieving user info: ", err)
		return c.String(http.StatusBadRequest, "Error retrieving user info")
	}

	user, err := database.FindApiUser(ctx, userInfo.ID)
	if err != nil {
		log.Println("User not found", err)
		return c.String(http.StatusNotFound, "User not found")
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Println("Error creating response:  ", err)
		return c.String(http.StatusInternalServerError, "Error creating response")
	}

	return c.String(http.StatusOK, string(response))
}

func getGuilds(c echo.Context) error {
	afterId := c.QueryParam("afterId")
	token := c.Request().Header.Get("Authorization")

	client, err := discordgo.New("Bearer " + token)
	if err != nil {
		log.Println("Error creating discord client: ", err)
		return c.String(http.StatusBadRequest, "Error creating discord client")
	}

	guilds, err := client.UserGuilds(100, "", afterId)

	if err != nil {
		log.Println("Error Fetching Guilds")
		return c.String(http.StatusBadRequest, "Error Fetching Guilds")
	}

	response, err := json.Marshal(guilds)
	if err != nil {
		log.Println("Error creating response:  ", err)
		return c.String(http.StatusInternalServerError, "Error creating response")
	}

	return c.String(http.StatusOK, string(response))
}

func getGuild(c echo.Context) error {
	guildId := c.Param("guildId")
	token := c.Request().Header.Get("Authorization")

	client, err := discordgo.New("Bearer " + token)
	if err != nil {
		log.Println("Error creating discord client: ", err)
		return c.String(http.StatusBadRequest, "Error creating discord client")
	}

	guilds, err := client.UserGuilds(100, "", "")

	var isMember = false
	for _, guild := range guilds {
		if guild.ID == guildId {
			isMember = true
			break
		}
	}

	if !isMember {
		log.Println("User not cannot access guild: ", guildId)
		return c.String(http.StatusBadRequest, "User not cannot access guild: "+guildId)
	}

	guild, err := database.FindProcessedGuild(ctx, guildId)

	if err != nil {
		log.Println("Error Fetching Guild", err)
		return c.String(http.StatusBadRequest, "Error Fetching Guild")
	}

	response, err := json.Marshal(guild)
	if err != nil {
		log.Println("Error creating response:  ", err)
		return c.String(http.StatusInternalServerError, "Error creating response")
	}

	return c.String(http.StatusOK, string(response))
}

// API UTILS

func retrieveCredentialsFromCode(code string) (apiModel.DiscordAuthResponse, error) {

	var responseJson apiModel.DiscordAuthResponse

	form := url.Values{}
	form.Add("client_id", config.ClientId)
	form.Add("client_secret", config.ClientSecret)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", config.OAuthRedirectURI)

	req, err := http.NewRequest("POST", "https://discord.com/api/v10/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Error creating auth request: ", err)
		return responseJson, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error making HTTP request: ", err)
		return responseJson, err
	}

	read, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Error reading auth response body: ", err)
		return responseJson, err
	}

	if res.StatusCode >= 400 {
		return responseJson, errors.New("Error retrieving code from discord: " + string(read))
	}

	err = json.Unmarshal(read, &responseJson)

	if err != nil {
		log.Println("Error UnMarshalling response body", err)
		return responseJson, err
	}

	return responseJson, nil
}

func getUserInfoFromAccessToken(token string) (*discordgo.User, error) {

	var user *discordgo.User
	client, err := discordgo.New("Bearer " + token)
	if err != nil {
		return user, err
	}
	user, err = client.User("@me")

	if err != nil {
		return user, err
	}
	return user, nil
}
