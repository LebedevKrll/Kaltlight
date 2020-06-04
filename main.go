package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	db      *sql.DB
)

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func listFromSql(values *sql.Rows, err error) []string {

	if err != nil {
		log.Fatal(err)
	}

	list := []string{}
	var id string
	for values.Next() {
		if err = values.Scan(&id); err != nil {
			log.Fatal(err)
		}
		list = append(list, id)
	}

	return list
}

func handlerRegister(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	rand.Seed(time.Now().UnixNano())
	token := randSeq(25)
	names, err := db.Query("SELECT name FROM users")
	listOfTheNames := listFromSql(names, err)

	for index, _ := range listOfTheNames {
		if name == listOfTheNames[index] {
			w.Write([]byte("This name is taken"))
			return
		}
	}

	db.Exec("INSERT INTO users (name, token) VALUES (?, ?);", name, token)
	w.Write([]byte("Welcome, " + name + "!" + " Your token is " + token))
}

func handlerCreate(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")
	text := r.URL.Query().Get("text")

	tokens, err := db.Query("SELECT token FROM users")
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {
			db.Exec("INSERT INTO texts (token, title, text) VALUES (?, ?, ?);", token, title, text)
			w.Write([]byte("File was created"))
			return
		}
	}
	w.Write([]byte("Invalid token"))
	return
}

func handlerShowFiles(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	tokens, err := db.Query("SELECT token FROM users")
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			listOfTheTitles := listFromSql(titles, err)

			w.Write([]byte("These are your files: "))
			for _, title := range listOfTheTitles {
				w.Write([]byte(title + " "))
			}
			return
		}
	}
	w.Write([]byte("Invalid token"))
	return

}

func handlerShowText(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	tokens, err := db.Query("SELECT token FROM users")
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			listOfTheTitles := listFromSql(titles, err)

			for index, _ = range listOfTheTitles {
				if title == listOfTheTitles[index] {

					texts, err := db.Query("SELECT text FROM texts WHERE texts.title = ? AND texts.token = ?;", title, token)
					listOfTheTexts := listFromSql(texts, err)

					w.Write([]byte("Text of the file is: "))
					for _, value := range listOfTheTexts {
						w.Write([]byte(value))
					}
					return
				}
			}
			w.Write([]byte("Invalid title"))
			return
		}
	}
	w.Write([]byte("Invalid token."))
	return

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	tokens, err := db.Query("SELECT token FROM users")
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			listOfTheTitles := listFromSql(titles, err)

			for index, _ = range listOfTheTitles {
				if title == listOfTheTitles[index] {
					db.Exec("DELETE FROM texts WHERE texts.title = ? AND texts.token = ?;", title, token)
					w.Write([]byte("File was deleted."))
					return
				}
			}
			w.Write([]byte("Invalid title."))
			return
		}
	}
	w.Write([]byte("Invalid token."))
	return

}

func main() {
	db, _ = sql.Open("sqlite3", "./data.db")
	db.Exec("CREATE TABLE IF NOT EXISTS users (name TEXT, token TEXT);")
	db.Exec("CREATE TABLE IF NOT EXISTS texts (token TEXT, title TEXT, text TEXT);")
	http.HandleFunc("/", handlerRegister)
	http.HandleFunc("/create", handlerCreate)
	http.HandleFunc("/showfiles", handlerShowFiles)
	http.HandleFunc("/showtext", handlerShowText)
	http.HandleFunc("/delete", handlerDelete)
	http.ListenAndServe(":8000", nil)
}
