package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Database string `json:"database"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Data struct {
	db     *sql.DB
	config Config
}

func (d *Data) init(c Config) {
	d.config = c
}

func (d *Data) start() {
	db, err := sql.Open("mysql", d.config.Username+":"+d.config.Password+"@/"+d.config.Database)

	d.db = db

	if err != nil {
		panic(err.Error())
	}
}

func (d *Data) stop() {
	d.db.Close()
}

func (d *Data) ping() {
	err := d.db.Ping()

	if err != nil {
		panic(err.Error())
	}

	log.Println("running...")
}

func (d *Data) insert(query string, values string) int64 {	
	t := time.Now()
	t.Format("2006-01-02 15:04:05")

	result, err := d.db.Exec(query, values, t)

	if err != nil {
		log.Fatal("insert", result, err)
	}

	result_id, err := result.LastInsertId()

	if err != nil {
		log.Fatal("insert result", err)
	}

	return result_id
}

func (d *Data) addCard(value string) int64 {
	return d.insert(`INSERT INTO cards(data, created_at) VALUES(?,?);`, value)
}

func (d *Data) addTag(value string) int64 {
	return d.insert(`INSERT INTO tags(data, created_at) VALUES(?,?);`, value)
}

func (d *Data) addCardTag(card_id, tag_id int64) int64 {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")

	result, err := d.db.Exec(`INSERT INTO cards_tags(card_id, tag_id, created_at) VALUES(?,?,?);`, card_id, tag_id, t)

	if err != nil {
		log.Fatal("addCardTag", result, err)
	}

	result_id, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	return result_id
}

func (d *Data) allRows(query string) {
	rows, err := d.db.Query(query)
	defer rows.Close()

	columns, err := rows.Columns()

	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		if err != nil {
			panic(err.Error())
		}

		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			fmt.Println(columns[i], ": ", value)
		}

		fmt.Println("-----------------------------------")
	}
}

func (d *Data) listCards() {
	d.allRows(`SELECT * FROM cards;`)
}

func (d *Data) listTags() {
	d.allRows(`SELECT * FROM tags;`)
}

func (d *Data) listCardTags() {
	d.allRows(`SELECT * FROM cards_tags;`)
}

func (c *Config) init() {
	configRaw, err := ioutil.ReadFile("config.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configRaw, &c)

	if err != nil {
		panic(err)
	}
}

func main() {
	var config Config

	config.init()

	var d Data

	d.init(config)

	d.start()

	defer d.stop()

	d.ping()

	values := os.Args[1:]

	for _, value := range values {
		card_id := d.addCard(value)
		tag_id := d.addTag(value)
		d.addCardTag(card_id, tag_id)
	}

	d.listCards()
	d.listTags()
	d.listCardTags()
}
