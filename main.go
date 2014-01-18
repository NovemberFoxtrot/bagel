package main

import (
	"database/sql"
	"encoding/json"
	"flag"
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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	check(err)
}

func (d *Data) stop() {
	d.db.Close()
}

func (d *Data) ping() {
	err := d.db.Ping()

	check(err)

	log.Println("running...")
}

func (d *Data) insert(query string, values... interface{}) int64 {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")

	values = append(values, t)

	result, err := d.db.Exec(query, values...)

	check(err)

	result_id, err := result.LastInsertId()

	check(err)

	return result_id
}

func (d *Data) addCard(answer, question, explanation string) int64 {
	return d.insert(`INSERT INTO cards(question, answer, explanation, created_at) VALUES(?,?,?,?);`, answer, question, explanation)
}

func (d *Data) addTag(value string) int64 {
	return d.insert(`INSERT INTO tags(data, created_at) VALUES(?,?);`, value)
}

func (d *Data) addCardTag(card_id, tag_id int64) int64 {
	return d.insert(`INSERT INTO cards_tags(card_id, tag_id, created_at) VALUES(?,?,?);`, card_id, tag_id)
}

func (d *Data) allRows(query string) {
	rows, err := d.db.Query(query)

	defer rows.Close()

	columns, err := rows.Columns()

	check(err)

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		check(err)

		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			fmt.Print(columns[i], " : ", value, " | ")
		}

		fmt.Println("")
	}
}

func (d *Data) listCards() {
	d.allRows(`SELECT * FROM cards;`)
}

func (d *Data) listTags() {
	d.allRows(`SELECT * FROM tags;`)
}

func (d *Data) listCardTags() {
	// d.allRows(`SELECT * FROM cards_tags;`)
	d.allRows(`SELECT 
	c.id,
	c.question,
	c.answer,
	c.explanation,
	c.created_at,
	t.data,
	t.created_at
	FROM cards c 
	JOIN cards_tags ct 
	ON c.id = ct.card_id 
	JOIN tags t 
	ON t.id = ct.tag_id`)
}

func (c *Config) init() {
	configRaw, err := ioutil.ReadFile("config.json")

	check(err)

	err = json.Unmarshal(configRaw, &c)

	check(err)
}

func main() {
	flag.Parse()

	var config Config

	config.init()

	var d Data

	d.init(config)

	d.start()

	defer d.stop()

	d.ping()

	values := os.Args[1:]

	if len(values) > 0 {
		card_id := d.addCard(values[0], values[1], values[2])
		tag_id := d.addTag(values[3])
		d.addCardTag(card_id, tag_id)
	}

	d.listCards()
	d.listTags()
	d.listCardTags()
}
