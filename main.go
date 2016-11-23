package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Database struct {
	Config    *Config
	Data      map[string]string
	Persister *time.Ticker
}

func NewDB() *Database {
	return &Database{
		Data:      nil,
		Config:    nil,
		Persister: nil,
	}
}

func (db *Database) Start() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db.Config = config
	db.Persister = time.NewTicker(time.Second * time.Duration(config.SaveFreq))

	db.LoadData()

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

	SocketStart(db.Config.SocketAddr)
	WebStart(db.Config.HttpAddr)
}

// Saving and loading

func (db *Database) Save() error {
	data, err := json.Marshal(db.Data)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(db.Config.DBFilepath, data, 0644)

	return err
}

func (db *Database) LoadData() {
	log.Println("Attempting to load", db.Config.DBFilepath)

	content, err := ioutil.ReadFile(db.Config.DBFilepath)

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
