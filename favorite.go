package main

import (
	"database/sql"
	"fmt"
	"time"
)

// CREATE TABLE favorites (
//   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   user_id INT,
//   status_id VARCHAR(64),
//   text VARCHAR(255),
//   screen_name VARCHAR(255),
//   created_at TIMESTAMP NOT NULL
// );

// Favorite by a Twitter user
type Favorite struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	StatusID   string    `json:"status_id"`
	Text       string    `json:"text"`
	ScreenName string    `json:"screen_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func favoriteTweetsPeriodically(db *sql.DB) {
	for true {
		fmt.Println("Favoriting Tweets")
		time.Sleep(5 * time.Hour)

		users := findAllUsers(db)
		for _, user := range users {
			keywords := findAllKeywordsByUser(user)
			for _, keyword := range keywords {
				search := searchForTweets(user, keyword)
				favoriteTweets(user, search.Statuses)
			}
		}
	}
}

func createFavorite(db *sql.DB, user *User, status *Status) {
	writeStmt, _ := db.Prepare("INSERT favorites SET user_id=?,status_id=?,text=?,screen_name=?")
	_, err := writeStmt.Exec(user.ID, fmt.Sprintf("%d", status.ID), status.Text, status.User.ScreenName)
	checkErr(err)
}

func alreadyFavorited(db *sql.DB, user *User, status *Status) bool {
	rows, err := db.Query("SELECT COUNT(*) AS matching FROM favorites WHERE user_id=? AND status_id=?", user.ID, fmt.Sprintf("%d", status.ID))
	checkErr(err)
	if rows.Next() {
		var matchingCount int
		rows.Scan(&matchingCount)
		return matchingCount > 0
	}
	return false
}

func findAllFavoritesByUser(db *sql.DB, user *User) []*Favorite {
	var favorites []*Favorite
	rows, err := db.Query("SELECT * FROM favorites WHERE user_id=?", user.ID)
	checkErr(err)
	for rows.Next() {
		f := new(Favorite)
		rows.Scan(&f.ID, &f.UserID, &f.StatusID, &f.Text, &f.ScreenName, &f.CreatedAt)
		favorites = append(favorites, f)
	}
	return favorites
}
