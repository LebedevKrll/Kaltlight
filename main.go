package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
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

func checkError(err error) {
	if err != nil {
		if _, osErr := os.Stat("./log.txt"); osErr == nil {
			f, _ := os.OpenFile("./log.txt", os.O_APPEND|os.O_WRONLY, 0600)
			f.WriteString(time.Now().String()[:19] + " " + err.Error() + "\n")
			f.Close()
			log.Fatal(err)
		} else {
			newFile, _ := os.Create("./log.txt")
			newFile.WriteString(time.Now().String()[:19] + " " + err.Error() + "\n")
			newFile.Close()
			log.Fatal(err)
		}
	}
	return
}

func listFromSql(values *sql.Rows, err error) []string {
	checkError(err)

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
	checkError(err)
	listOfTheNames := listFromSql(names, err)

	for index, _ := range listOfTheNames {
		if name == listOfTheNames[index] {
			jsonText, jsonErr := json.Marshal("This name is taken")
			checkError(jsonErr)
			w.Write(jsonText[1 : len(jsonText)-1])
			return
		}
	}

	db.Exec("INSERT INTO users (name, token) VALUES (?, ?);", name, token)
	jsonText, jsonErr := json.Marshal("Welcome, " + name + "!" + " Your token is " + token)
	checkError(jsonErr)
	w.Write(jsonText[1 : len(jsonText)-1])
}

func handlerCreate(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")
	text := r.URL.Query().Get("text")

	tokens, err := db.Query("SELECT token FROM users")
	checkError(err)
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {
			db.Exec("INSERT INTO texts (token, title, text) VALUES (?, ?, ?);", token, title, text)
			jsonText, jsonErr := json.Marshal("File was created")
			checkError(jsonErr)
			w.Write(jsonText[1 : len(jsonText)-1])
			return
		}
	}
	jsonText, jsonErr := json.Marshal("Invalid token")
	checkError(jsonErr)
	w.Write(jsonText[1 : len(jsonText)-1])
	return
}

func handlerShowFiles(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	tokens, err := db.Query("SELECT token FROM users")
	checkError(err)
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			checkError(err)
			listOfTheTitles := listFromSql(titles, err)

			jsonText, jsonErr := json.Marshal("These are your files: ")
			checkError(jsonErr)
			w.Write(jsonText[1 : len(jsonText)-1])
			for _, title := range listOfTheTitles {
				jsonText, jsonErr := json.Marshal(title + " ")
				checkError(jsonErr)
				w.Write(jsonText[1 : len(jsonText)-1])
			}
			return
		}
	}
	jsonText, jsonErr := json.Marshal("Invalid token")
	checkError(jsonErr)
	w.Write(jsonText[1 : len(jsonText)-1])
	return

}

func handlerShowText(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	tokens, err := db.Query("SELECT token FROM users")
	checkError(err)
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			checkError(err)
			listOfTheTitles := listFromSql(titles, err)

			for index, _ = range listOfTheTitles {
				if title == listOfTheTitles[index] {

					texts, err := db.Query("SELECT text FROM texts WHERE texts.title = ? AND texts.token = ?;", title, token)
					checkError(err)
					listOfTheTexts := listFromSql(texts, err)

					jsonText, jsonErr := json.Marshal("Text of the file is: ")
					checkError(jsonErr)
					w.Write(jsonText[1 : len(jsonText)-1])
					for _, text := range listOfTheTexts {
						jsonText, jsonErr := json.Marshal(text)
						checkError(jsonErr)
						w.Write(jsonText[1 : len(jsonText)-1])
					}
					return
				}
			}
			jsonText, jsonErr := json.Marshal("Invalid title")
			checkError(jsonErr)
			w.Write(jsonText[1 : len(jsonText)-1])
			return
		}
	}
	jsonText, jsonErr := json.Marshal("Invalid token")
	checkError(jsonErr)
	w.Write(jsonText[1 : len(jsonText)-1])
	return

}

func handlerDelete(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	title := r.URL.Query().Get("title")

	tokens, err := db.Query("SELECT token FROM users")
	checkError(err)
	listOfTheTokens := listFromSql(tokens, err)

	for index, _ := range listOfTheTokens {
		if token == listOfTheTokens[index] {

			titles, err := db.Query("SELECT title FROM texts WHERE texts.token = ?;", token)
			checkError(err)
			listOfTheTitles := listFromSql(titles, err)

			for index, _ = range listOfTheTitles {
				if title == listOfTheTitles[index] {
					db.Exec("DELETE FROM texts WHERE texts.title = ? AND texts.token = ?;", title, token)
					jsonText, jsonErr := json.Marshal("File was deleted")
					checkError(jsonErr)
					w.Write(jsonText[1 : len(jsonText)-1])
					return
				}
			}
			jsonText, jsonErr := json.Marshal("Invalid title")
			checkError(jsonErr)
			w.Write(jsonText[1 : len(jsonText)-1])
			return
		}
	}
	jsonText, jsonErr := json.Marshal("Invalid token")
	checkError(jsonErr)
	w.Write(jsonText[1 : len(jsonText)-1])
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
