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
)

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func listfromsql(values *sql.Rows, err error) []string {

	if err != nil {
		log.Fatal(err)
	}

	llist := []string{}
	var id string
	for values.Next() {
		if err = values.Scan(&id); err != nil {
			log.Fatal(err)
		}
		llist = append(llist, id)
	}

	return llist
}

func handler_reg(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	rand.Seed(time.Now().UnixNano())
	token := randSeq(25)

	db, _ := sql.Open("sqlite3", "./data.db")
	names, err := db.Query("SELECT name FROM users")
	llist := listfromsql(names, err)

	for index, _ := range llist {
		if name == llist[index] {
			w.Write([]byte("This name is taken"))
			return
		}
	}

	db.Exec("INSERT INTO users (name, token) VALUES (?, ?);", name, token)
	w.Write([]byte("Welcome, " + name + "!" + " Your token is " + token))
}

func handler_create(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")
	text := r.URL.Query().Get("text")

	db, _ := sql.Open("sqlite3", "./data.db")
	tokens, err := db.Query("SELECT token FROM users")
	llist := listfromsql(tokens, err)

	for index, _ := range llist {
		if token == llist[index] {
			db.Exec("INSERT INTO texts (token, title, text) VALUES (?, ?, ?);", token, title, text)
			w.Write([]byte("File was created"))
			return
		}
	}
	w.Write([]byte("Invalid token"))
	return
}

func handler_catall(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	db, _ := sql.Open("sqlite3", "./data.db")
	tokens, err := db.Query("SELECT token FROM users")
	llist := listfromsql(tokens, err)

	for index, _ := range llist {
		if token == llist[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			llist2 := listfromsql(titles, err)

			w.Write([]byte("These are your files: "))
			for _, title := range llist2 {
				w.Write([]byte(title + " "))
			}
			return
		}
	}
	w.Write([]byte("Invalid token"))
	return

}

func handler_cat(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	db, _ := sql.Open("sqlite3", "./data.db")
	tokens, err := db.Query("SELECT token FROM users")
	llist := listfromsql(tokens, err)

	for index, _ := range llist {
		if token == llist[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			llist2 := listfromsql(titles, err)

			for index, _ = range llist2 {
				if title == llist2[index] {

					texts, err := db.Query("SELECT text FROM texts WHERE texts.title = ? AND texts.token = ?;", title, token)
					llist3 := listfromsql(texts, err)

					w.Write([]byte("Text of the file is: "))
					for _, value := range llist3 {
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

func handler_del(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	db, _ := sql.Open("sqlite3", "./data.db")
	tokens, err := db.Query("SELECT token FROM users")
	llist := listfromsql(tokens, err)

	for index, _ := range llist {
		if token == llist[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			llist2 := listfromsql(titles, err)

			for index, _ = range llist2 {
				if title == llist2[index] {
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
	db, _ := sql.Open("sqlite3", "./data.db")
	db.Exec("CREATE TABLE users (name TEXT, token TEXT);")
	db.Exec("CREATE TABLE texts (token TEXT, title TEXT, text TEXT);")
	http.HandleFunc("/", handler_reg)
	http.HandleFunc("/create", handler_create)
	http.HandleFunc("/catall", handler_catall)
	http.HandleFunc("/cat", handler_cat)
	http.HandleFunc("/del", handler_del)
	http.ListenAndServe(":8000", nil)
}
