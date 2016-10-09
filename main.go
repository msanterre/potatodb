package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	listenAddress = "127.0.0.1:5050"
	dbFilepath    = "data.json"
)

type Database struct {
	Data      map[string]string
	Persister *time.Ticker
}

func NewDB() *Database {
	return &Database{
		Data:      nil,
		Persister: time.NewTicker(time.Second * 10),
	}
}

func (db *Database) Start() {
	db.Load()

	go func() {
		for range db.Persister.C {
			log.Println("Saving...")
			err := db.Save()

			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	defer db.Persister.Stop()

	http.HandleFunc("/db/", DBMethod)

	log.Printf("Listening on %s\n", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// Saving and loading

func (db *Database) Save() error {
	data, err := json.Marshal(db.Data)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(dbFilepath, data, 0644)

	return err
}

func (db *Database) Load() {
	log.Println("Attempting to load", dbFilepath)

	content, err := ioutil.ReadFile(dbFilepath)

	if err != nil {
		log.Println("Could not read file. Starting from scratch")
		db.Data = make(map[string]string)
		return
	}

	var data map[string]string

	err = json.Unmarshal(content, &data)

	if err != nil {
		log.Println("Could not decode file. Starting from scratch")
		db.Data = make(map[string]string)
		return
	}

	log.Println("Finished loading")

	db.Data = data
}

// API

func Error(w http.ResponseWriter, req *http.Request, err string) {
	io.WriteString(w, err)
}

func DBMethod(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.RequestURI, "/db/")

	if key == "" {
		NotFoundError(w)
		return
	}

	switch req.Method {
	case "GET":
		if value, ok := DB.Data[key]; ok {
			io.WriteString(w, value)
		} else {
			NotFoundError(w)
		}
	case "POST", "PUT", "PATCH":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ServerError(w, err)
			log.Println(err)
			return
		}
		DB.Data[key] = string(body)

		w.WriteHeader(200)
	case "DELETE":
		delete(DB.Data, key)
		w.WriteHeader(200)
	}
}

func NotFoundError(w http.ResponseWriter) {
	w.WriteHeader(404)
	io.WriteString(w, "Not found")
}

func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	io.WriteString(w, err.Error())
}

var DB *Database

func main() {
	DB = NewDB()
	DB.Start()
}
