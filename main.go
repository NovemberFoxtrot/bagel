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

func check(err error) {
	if err != nil {
		log.Fatal("error:", err)
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

func (d *Data) add(stack, question, answer string) int64 {
	return d.insert(`INSERT INTO cards(stack, question, answer, created_at) VALUES(?,?,?,?);`, stack, question, answer)
}

func (d *Data) allRows(query string) []string {
	var results []string

	rows, err := d.db.Query(query)

	defer rows.Close()

	columns, err := rows.Columns()

	check(err)

	count := 0

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		count += 1
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
			results = append(results, value)
		}

		fmt.Println("")
	}

	return results
}

func (c *Config) init() {
	configRaw, err := ioutil.ReadFile("config.json")

	check(err)

	err = json.Unmarshal(configRaw, &c)

	check(err)
}

func (d *Data) next(stack string) []string {
	return d.allRows(`SELECT id, question, answer, ROUND(correct / (incorrect + 1)) AS card_status FROM cards WHERE stack = '` + stack + `' ORDER BY card_status, RAND() LIMIT 1;`)
}

func (d *Data) correct(id string) {
	_, err := d.db.Exec("UPDATE cards SET `correct` = `correct` + 1 WHERE id = ?", id)
	check(err)
}

func (d *Data) incorrect(id string) {
	_, err := d.db.Exec("UPDATE cards SET `incorrect` = `incorrect` + 1 WHERE id = ?", id)
	check(err)
}

func (d *Data) learn(stack string) {
	for {
		current := d.next(stack)

		id := current[0]

		var response string
		fmt.Printf("%s %s", id, current[1])
		fmt.Scanf("%s", &response)

		fmt.Printf("%s %s %s", id, current[2], "(y/N)? :")
		fmt.Scanf("%s", &response)

		switch response {
		case "", "n", "N":
			d.incorrect(id)
		default:
			d.correct(id)
		}
	}
}

func usage() {
	fmt.Printf("usage: %s <add/learn> <options>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	var config Config
	config.init()

	var d Data
	d.init(config)
	d.start()
	defer d.stop()
	d.ping()

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 5 {
			usage()
		}

		d.add(os.Args[2], os.Args[3], os.Args[4])
		d.add(os.Args[2], os.Args[4], os.Args[3])
	case "learn":
		if len(os.Args) != 3 {
			usage()
		}

		d.learn(os.Args[2])
	default:
		usage()
	}
}
