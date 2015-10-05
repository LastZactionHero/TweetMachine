package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

import "net/http"

import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func main() {
	fmt.Println("Tweet Machine")
	godotenv.Load()
	db = connectToDatabase()

	go favoriteTweetsPeriodically(db)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login_callback", loginCallbackHandler)
	http.HandleFunc("/keywords", keywordHandler)
	http.HandleFunc("/favorites", favoritesHandler)
	http.HandleFunc("/", defaultHandler)

	http.ListenAndServe(":9000", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	var keywordMatcher = regexp.MustCompile(`\/keywords\/[0-9]+`)

	switch {
	case r.URL.String() == "/" && r.Method == "GET":
		loginHandler(w, r)
	case keywordMatcher.MatchString(r.URL.String()):
		keywordHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	oauthVerifier := r.URL.Query()["oauth_verifier"][0]

	consumer := newOauthConsumer()

	secretCookie, _ := r.Cookie("request_token_secret")
	tokenCookie, _ := r.Cookie("request_token_token")

	secret := secretCookie.Value
	token := tokenCookie.Value

	requestToken := restoreOauthRequestToken(token, secret)

	accessToken, _ := consumer.AuthorizeToken(&requestToken, oauthVerifier)
	storeOAuthAccessToken(accessToken)
	sessionStoreUser(w, accessToken.AdditionalData["screen_name"])

	fmt.Fprintf(w, fmt.Sprintf("Access Token Received: %s", accessToken.AdditionalData["screen_name"]))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	consumer := newOauthConsumer()

	callbackURL := fmt.Sprintf("%s/login_callback", os.Getenv("API_HOST"))
	requestToken, loginURL, _ := consumer.GetRequestTokenAndUrl(callbackURL)

	secretCookie := http.Cookie{Name: "request_token_secret", Value: requestToken.Secret}
	tokenCookie := http.Cookie{Name: "request_token_token", Value: requestToken.Token}

	http.SetCookie(w, &secretCookie)
	http.SetCookie(w, &tokenCookie)

	http.Redirect(w, r, loginURL, http.StatusFound)
}

func connectToDatabase() *sql.DB {
	url := fmt.Sprintf("%s@/%s?charset=utf8&parseTime=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME"))
	db, err := sql.Open("mysql", url)
	err = db.Ping()
	checkErr(err)
	return db
}

func keywordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := sessionRestoreUser(r, db)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == "GET":
		keywords := findAllKeywordsByUser(user)

		var arr []KeywordJSON
		for _, keyword := range keywords {
			arr = append(arr, KeywordJSON{ID: keyword.ID, Keyword: keyword.Keyword})
		}

		data, _ := json.Marshal(arr)
		w.Write(data)
	case r.Method == "POST":
		keywordStr := r.FormValue("keyword")
		keyword := Keyword{Keyword: keywordStr, UserID: user.ID}
		err := keyword.Store(db)
		if err != nil {
			respondWithError(w, err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		}
	case r.Method == "DELETE":
		var idMatcher = regexp.MustCompile(`[0-9]+`)
		idStr := idMatcher.FindString(r.URL.String())
		id, _ := strconv.Atoi(idStr)
		keyword, err := findKeywordByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		keyword.Delete(db)

	default:
		panic("Invalid Method")
	}
}

func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := sessionRestoreUser(r, db)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	favorites := findAllFavoritesByUser(db, user)
	data, _ := json.Marshal(favorites)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, err error) {
	m := make(map[string]string)
	m["error"] = err.Error()
	error, _ := json.Marshal(m)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(error)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
