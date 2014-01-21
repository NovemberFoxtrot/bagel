package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func (d *Data) insert(query string, values ...interface{}) int64 {
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
	d.allRows(`SELECT 
	c.id,
	c.question,
	c.answer,
	c.explanation,
	c.created_at,
	t.data,
	t.created_at
	FROM cards AS c 
	JOIN cards_tags AS ct 
	ON c.id = ct.card_id 
	JOIN tags AS t 
	ON t.id = ct.tag_id`)
}

func (c *Config) init() {
	configRaw, err := ioutil.ReadFile("config.json")

	check(err)

	err = json.Unmarshal(configRaw, &c)

	check(err)
}

func (d *Data) parseCSV(theFile string) {
	csvFile, err := os.Open(theFile)
	defer csvFile.Close()
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(csvFile)

	for {
		fields, err := csvReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		data := fields[0]
		// cost1 := fields[1]
		// cost2 := fields[2]
		// cost3 := fields[3]
		speech := fields[4]
		// something := fields[5]
		// form := fields[6]
		// dictionary := fields[7]
		hiragana := fields[8]
		notes := fields[9]

		// fmt.Println(data, cost1, cost2, cost3, speech, something, form, dictionary, hiragana, notes)

		card_id := d.addCard(data, hiragana, notes)
		tag_id := d.addTag(speech)

		d.addCardTag(card_id, tag_id)
	}
}

func main() {
	var theFile string

	flag.StringVar(&theFile, "f", "file", "")

	flag.Parse()

	var config Config

	config.init()

	var d Data

	d.init(config)

	d.start()

	defer d.stop()

	d.ping()

	d.parseCSV(theFile)

	//d.listCards()
	//d.listTags()
	d.listCardTags()
}
