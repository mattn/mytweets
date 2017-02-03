package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

var (
	addr = flag.String("a", ":8989", "server address")
)

type tweet struct {
	TweetID   string `json:"tweet_id"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func main() {
	flag.Parse()

	os.Remove("tweets.db")

	b, err := exec.Command("sqlite3", "-separator", ",", "tweets.db", ".import tweets.csv tweets").CombinedOutput()
	if err != nil {
		log.Fatalf("csv import error: %v: %s", err, string(b))
	}

	conn, err := sql.Open("sqlite3", "tweets.db")
	if err != nil {
		log.Fatalf("connect database error: %v", err)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, req *http.Request) {
		q := req.FormValue("q")

		query := `
		select tweet_id, text, timestamp from tweets where text like $1
		`
		if req.FormValue("nort") != "" {
			query += " and retweeted_status_user_id == ''"
		}
		query += " order by timestamp desc"
		rows, err := conn.Query(query, fmt.Sprintf("%%%s%%", q))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Write([]byte("["))
		first := true
		enc := json.NewEncoder(w)
		for rows.Next() {
			if first {
				first = false
			} else {
				w.Write([]byte(","))
			}
			var t tweet
			err = rows.Scan(&t.TweetID, &t.Text, &t.Timestamp)
			if err != nil {
				break
			}
			err = enc.Encode(&t)
			if err != nil {
				break
			}
		}
		w.Write([]byte("]"))
	})
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(*addr, nil)
}
