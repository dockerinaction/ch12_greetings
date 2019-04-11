package main

import (
	"database/sql"
	"fmt"
	"strings"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	hostname string
)

func main() {
	var err error
	
	hostname, err = os.Hostname()
	log.Println("Initializing api server on", hostname)
	
	postgres_password_file := os.Getenv("POSTGRES_PASSWORD_FILE")
	log.Println(os.ExpandEnv("Will read postgres password from '${POSTGRES_PASSWORD_FILE}'"))

	postgres_password, err := ioutil.ReadFile(postgres_password_file)
	if err != nil {
		log.Fatal(err)
	}

	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		strings.TrimSpace(string(postgres_password)),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)

	if os.Getenv("DEBUG") == "true" {
		log.Println(connInfo)
	}

	db, err = sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = db.Ping(); err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`create table if not exists counter (
			id serial primary key,
			val integer not null
		)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/counter", serveCounter)

	log.Println("Initialization complete, starting http service on", hostname)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func serveIndex(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.Header().Set("Connection", "close")
	
	fmt.Fprintln(resp, "Welcome to the API Server!")
	fmt.Fprintln(resp, fmt.Sprintf("Container with id %s responded at %s", hostname, time.Now().UTC()))
}

func serveCounter(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.Header().Set("Connection", "close")

	fmt.Fprintln(resp, "Container ID:", hostname)
	fmt.Fprintln(resp, "POSTGRES_HOST:", os.Getenv("POSTGRES_HOST"))
	fmt.Fprintln(resp, "POSTGRES_PORT:", os.Getenv("POSTGRES_PORT"))

	_, err := db.Exec("insert into counter(val) values(0)")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id from counter")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(resp, "ID: %d\n", id)
	}
}
