package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Database string `json:"database"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func main() {
	configRaw, err := ioutil.ReadFile("config.json")

	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(configRaw, &config)

	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", config.Username+":"+config.Password+"@/"+config.Database)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	log.Println("running...")
}
