package apiModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscordAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type UserAuth struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserGuildInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Owner        bool   `json:"owner"`
	IsBotPresent bool   `json:"is_bot_present"`
}

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	UserId        string             `json:"user_id" bson:"user_id"`
	UserName      string             `json:"user_name" bson:"user_name"`
	AccentColor   int                `json:"accent_color" bson:"accent_color"`
	Discriminator string             `json:"discriminator" bson:"discriminator"`
	Verified      bool               `json:"verified" bson:"verified"`
	Locale        string             `json:"locale" bson:"locale"`
	PremiumType   int                `json:"premium_type" bson:"premium_type"`
	Banner        string             `json:"banner" bson:"banner"`
	Avatar        string             `json:"avatar" bson:"avatar"`
	MFAEnabled    bool               `json:"mfa_enabled" bson:"mfa_enabled"`
}
