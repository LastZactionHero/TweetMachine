package main

import (
	"database/sql"
	"net/http"
)

func sessionStoreUser(w http.ResponseWriter, username string) {
	usernameCookie := http.Cookie{Name: "username", Value: username}
	http.SetCookie(w, &usernameCookie)
}

func sessionRestoreUser(r *http.Request, db *sql.DB) *User {
	usernameCookie, _ := r.Cookie("username")
	if usernameCookie != nil {
		username := usernameCookie.Value

		rows, err := db.Query("SELECT * FROM users WHERE username=?", username)
		checkErr(err)

		if rows.Next() {
			var user User
			err = rows.Scan(&user.ID, &user.Token, &user.Secret, &user.Username)
			checkErr(err)
			return &user
		}
	}
	return nil
}
