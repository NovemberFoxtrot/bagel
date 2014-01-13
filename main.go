package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

func (d *Data) add(value string) {
	result, err := d.db.Exec(`INSERT INTO tags(data) VALUES(?);`, value)

	if err != nil {
		log.Fatal(result, err)
	}
}

func (d *Data) list() {
	rows, err := d.db.Query(`SELECT * FROM tags;`)
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
		d.add(value)
	}

	d.list()
}
