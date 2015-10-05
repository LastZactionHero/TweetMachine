package main

import (
	"os"

	"github.com/mrjones/oauth"
)

var gAccessToken *oauth.AccessToken

func newOauthConsumer() *oauth.Consumer {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")

	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	}
	return oauth.NewConsumer(consumerKey, consumerSecret, provider)
}

func storeOAuthAccessToken(accessToken *oauth.AccessToken) {
	gAccessToken = accessToken

	username := accessToken.AdditionalData["screen_name"]
	findOrCreateUser(username,
		accessToken.Token,
		accessToken.Secret)
}

func restoreOauthRequestToken(token string, secret string) oauth.RequestToken {
	t := oauth.RequestToken{Token: token, Secret: secret}
	return t
}
