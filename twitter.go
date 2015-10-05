package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mrjones/oauth"
)

// Search Twitter Search
type Search struct {
	Statuses []*Status `json:"statuses"`
}

type twitterUser struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// Status Twitter Status
type Status struct {
	Text              string      `json:"text"`
	ID                int         `json:"id"`
	User              twitterUser `json:"user"`
	InReplyToStatusID int         `json:"in_reply_to_status_id"`
}

func searchForTweets(user *User, keyword *Keyword) Search {
	accessToken := oauth.AccessToken{Token: user.Token, Secret: user.Secret}
	consumer := newOauthConsumer()

	searchParams := make(map[string]string)
	searchParams["q"] = keyword.Keyword
	resp, err := consumer.Get("https://api.twitter.com/1.1/search/tweets.json",
		searchParams,
		&accessToken)
	checkErr(err)
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	var search Search
	json.Unmarshal(bytes, &search)

	return search
}

func favoriteTweets(user *User, statuses []*Status) {
	accessToken := oauth.AccessToken{Token: user.Token, Secret: user.Secret}
	consumer := newOauthConsumer()

	for _, status := range statuses {
		if status.InReplyToStatusID > 0 {
			continue
		} else if alreadyFavorited(db, user, status) {
			continue
		} else {
			favoriteParams := make(map[string]string)
			favoriteParams["id"] = fmt.Sprintf("%d", status.ID)
			resp, err := consumer.Post("https://api.twitter.com/1.1/favorites/create.json",
				favoriteParams,
				&accessToken)
			if err == nil {
				createFavorite(db, user, status)
			}
			defer resp.Body.Close()
		}
	}
}
