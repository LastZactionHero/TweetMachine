package main

import "database/sql"

// CREATE TABLE users (
//   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   token VARCHAR(256),
//   secret VARCHAR(256),
//   username VARCHAR(64)
// );

// User object
type User struct {
	ID       int
	Username string
	Token    string
	Secret   string
}

func findOrCreateUser(username string, token string, secret string) {
	// Update exiting user
	rows, err := db.Query("SELECT id FROM users WHERE username=?", username)
	checkErr(err)

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		checkErr(err)

		updateStmt, _ := db.Prepare("update users set token=?,secret=? WHERE id=?")
		_, err := updateStmt.Exec(token, secret, id)
		checkErr(err)
	} else {
		writeStmt, _ := db.Prepare("INSERT users SET username=?,token=?,secret=?")
		_, err := writeStmt.Exec(db, username, token, secret)
		checkErr(err)
	}

}

func findAllUsers(db *sql.DB) []*User {
	var users []*User
	rows, err := db.Query("SELECT id,username,token,secret FROM users")
	checkErr(err)
	for rows.Next() {
		u := new(User)
		rows.Scan(&u.ID, &u.Username, &u.Token, &u.Secret)
		users = append(users, u)
	}
	return users
}
