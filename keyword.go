package main

import (
	"database/sql"
	"fmt"
)

// CREATE TABLE keywords (
// 	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
// 	user_id INT,
// 	keyword VARCHAR(256)
// );

// Keyword object
type Keyword struct {
	ID      int
	Keyword string
	UserID  int
}

// KeywordJSON object
type KeywordJSON struct {
	ID      int    `json:"id"`
	Keyword string `json:"keyword"`
}

// Store a keyword to database
func (k Keyword) Store(db *sql.DB) error {
	if len(k.Keyword) == 0 {
		return fmt.Errorf("Keyword must be present")
	}

	rows, err := db.Query("SELECT COUNT(*) AS matching FROM keywords WHERE keyword=? AND user_id=?", k.Keyword, k.UserID)
	if rows.Next() {
		var matchingCount int
		rows.Scan(&matchingCount)
		if matchingCount > 0 {
			return fmt.Errorf("Keyword must be unique")
		}
	}
	stmt, err := db.Prepare("INSERT keywords SET user_id=?,keyword=?")
	checkErr(err)

	_, err = stmt.Exec(k.UserID, k.Keyword)
	return err
}

// Delete a keyword from the database
func (k Keyword) Delete(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM keywords WHERE id=?")
	checkErr(err)
	_, err = stmt.Exec(k.ID)
	return err
}

func findKeywordByID(id int) (*Keyword, error) {
	rows, _ := db.Query("SELECT id,keyword FROM keywords WHERE id=?", id)
	if rows.Next() {
		var keyword = new(Keyword)
		rows.Scan(&keyword.ID, &keyword.Keyword)
		return keyword, nil
	}
	return nil, fmt.Errorf("Not found")
}

// List all user keywords
func findAllKeywordsByUser(user *User) []*Keyword {
	rows, err := db.Query("SELECT id,keyword FROM keywords WHERE user_id=?", user.ID)
	checkErr(err)

	var keywords []*Keyword
	keywords = make([]*Keyword, 0)

	for rows.Next() {
		keyword := new(Keyword)
		rows.Scan(&keyword.ID, &keyword.Keyword)
		keywords = append(keywords, keyword)
	}

	return keywords
}
