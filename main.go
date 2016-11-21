package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

const (
	webListenAddress    = "127.0.0.1:5050"
	socketListenAddress = "127.0.0.1:5051"
	dbFilepath          = "data.json"
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

	SocketStart(socketListenAddress)
	WebStart(webListenAddress)
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

var DB *Database

func main() {
	DB = NewDB()
	DB.Start()
}
