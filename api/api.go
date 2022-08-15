package api

import (
	"context"
	"dc-stats/api/api-model"
	"dc-stats/config"
	"dc-stats/database"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"

	"github.com/iris-contrib/middleware/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var appCtx context.Context
var apiLogger *golog.Logger

// @title       DC-STATS user API
// @version     0.0.3
// @description This api helps user see their data using discord code grant authentication
// @license.name MIT License
// @host     server.dc-stats.com
// @BasePath /v1
func Start(appContext context.Context) {

	appCtx = appContext
	irisApp := iris.New()

	// Logger
	irisApp.Logger().SetPrefix("[DC-STATS-API] ")
	irisApp.Logger().SetTimeFormat(time.RFC3339)
	apiLogger = irisApp.Logger()

	// CORS
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
		Debug:            true,
	})

	irisApp.UseRouter(corsConfig)

	// v1 path
	{
		// User Auth
		appApi := irisApp.Party("/v1")
		appApi.Post("/auth/:code", authentication)

		//UserAPI
		{
			userApi := appApi.Party("/user")

			userApi.Get("/guilds", getGuilds)

			userApi.Get("/guilds/:guildId", getGuild)

			userApi.Get("", getUser)
		}
	}

	for _, route := range irisApp.GetRoutes() {
		apiLogger.Info("Mapped Route: " + "[" + route.Method + "] " + route.Path)
	}

	apiLogger.Error(irisApp.Listen(":8080", iris.WithOptimizations))
}

// @id      authUser
// @Summary Authenticates a user by code grant
// @Tags    authentication
// @Accept  json
// @Produce json
// @Param   code path     string true "Discord code grant"
// @Success 200  {object} apiModel.UserAuth
// @Failure 400  {string} string
// @Router  /auth/{code} [post]
func authentication(irisCtx iris.Context) {
	code := irisCtx.Params().Get("code")

	authentication, err := retrieveCredentialsFromCode(code)
	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON("Error retrieving token")
		apiLogger.Error("Error retrieving token: ", err)
		return
	}

	accessToken := authentication.AccessToken

	userInfo, err := getUserInfoFromAccessToken(accessToken)

	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON("Error retrieving user info")
		apiLogger.Error("Error retrieving user info: ", err)
		return
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
	database.SaveOrUpdatApiUser(user, appCtx)

	userAuth := apiModel.UserAuth{
		UserId:       userInfo.ID,
		AccessToken:  authentication.AccessToken,
		RefreshToken: authentication.RefreshToken,
		ExpiresIn:    authentication.ExpiresIn}

	irisCtx.StatusCode(iris.StatusOK)
	irisCtx.JSON(userAuth)
}

// @id      getUser
// @Summary Gets a user given its authentication token
// @Tags    user
// @Accept  json
// @Produce json
// @Param   Authorization header   string true  "Discord authentication token"
// @Success 200  {object} apiModel.User
// @Failure 400  {string} string
// @Router  /user [get]
func getUser(irisCtx iris.Context) {
	token := irisCtx.Request().Header.Get("Authorization")

	userInfo, err := getUserInfoFromAccessToken(token)
	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON("Error retrieving user info")
		apiLogger.Error("Error retrieving user info: ", err)
		return
	}

	user, err := database.FindApiUser(appCtx, userInfo.ID)
	if err != nil {
		irisCtx.StatusCode(iris.StatusNotFound)
		irisCtx.JSON("User not found")
		apiLogger.Error("User not found", err)
		return
	}

	irisCtx.StatusCode(iris.StatusOK)
	irisCtx.JSON(user)
}

// @id      getGuilds
// @Summary Gets a user's guilds given its authentication token
// @Tags    guild
// @Accept  json
// @Produce json
// @Param   Authorization header   string true  "Discord authentication token"
// @Param   afterId       query    string false "Get guilds after this id"
// @Success 200           {array}  apiModel.UserGuildInfo
// @Failure 400           {string} string
// @Router  /user/guilds [get]
func getGuilds(irisCtx iris.Context) {
	afterId := irisCtx.URLParam("afterId")
	token := irisCtx.Request().Header.Get("Authorization")

	client, err := discordgo.New("Bearer " + token)
	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON("Error creating discord client")
		apiLogger.Error("Error creating discord client: ", err)
		return
	}

	guilds, err := client.UserGuilds(100, "", afterId)

	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON("Error Fetching Guilds")
		apiLogger.Error("Error Fetching Guilds: ", err)
		return
	}

	resultGuilds := []apiModel.UserGuildInfo{}
	for _, guild := range guilds {

		var found bool
		_, err := database.FindProcessedGuild(appCtx, guild.ID)
		if err != nil {
			found = false
		} else {
			found = true
		}
		resultGuilds = append(resultGuilds, apiModel.UserGuildInfo{
			ID:           guild.ID,
			Name:         guild.Name,
			Icon:         guild.Icon,
			Owner:        guild.Owner,
			IsBotPresent: found,
		})
	}

	irisCtx.StatusCode(iris.StatusOK)
	irisCtx.JSON(resultGuilds)

}

// @id      getGuild
// @Summary Gets a guild only if the user is inside it and we have a record of it
// @Tags    guild
// @Accept  json
// @Produce json
// @Param   Authorization header   string true "Discord authentication token"
// @Param   guildId       path     string true "Guild id"
// @Success 200           {object} processedModel.Guild
// @Failure 400           {string} string
// @Failure 404           {string} string
// @Failure 403           {string} string
// @Router  /user/guilds/{guildId} [get]
func getGuild(irisCtx iris.Context) {
	guildId := irisCtx.Params().Get("guildId")
	token := irisCtx.Request().Header.Get("Authorization")

	isMember, err := isUserInGuild(guildId, token)

	if err != nil {
		irisCtx.StatusCode(iris.StatusBadRequest)
		irisCtx.JSON(err.Error())
		apiLogger.Error(err)
	}

	if !isMember {
		irisCtx.StatusCode(iris.StatusForbidden)
		irisCtx.JSON("User not cannot access guild: " + guildId)
		apiLogger.Error("User not cannot access guild: ", guildId)
		return
	}

	guild, err := database.FindProcessedGuild(appCtx, guildId)

	if err != nil {
		irisCtx.StatusCode(iris.StatusNotFound)
		irisCtx.JSON("Guild not found")
		apiLogger.Error("Guild not found", err)
		return
	}

	irisCtx.StatusCode(iris.StatusOK)
	irisCtx.JSON(guild)
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
		apiLogger.Error("Error creating auth request: ", err)
		return responseJson, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		apiLogger.Error("Error making HTTP request: ", err)
		return responseJson, err
	}

	read, err := io.ReadAll(res.Body)

	if err != nil {
		apiLogger.Error("Error reading auth response body: ", err)
		return responseJson, err
	}

	if res.StatusCode >= 400 {
		return responseJson, errors.New("Error retrieving code from discord: " + string(read))
	}

	err = json.Unmarshal(read, &responseJson)

	if err != nil {
		apiLogger.Error("Error UnMarshalling response body", err)
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

func isUserInGuild(guildId string, token string) (bool, error) {

	client, err := discordgo.New("Bearer " + token)
	if err != nil {

		return false, errors.New("Error creating discord client")
	}
	afterId := ""
	guilds, err := client.UserGuilds(100, "", afterId)
	var isMember = false

	for len(guilds) > 0 && err == nil {

		for _, guild := range guilds {
			if guild.ID == guildId {
				isMember = true
				break
			}
		}
		afterId = guilds[len(guilds)-1].ID
		guilds, err = client.UserGuilds(100, "", afterId)
	}
	return isMember, nil
}
