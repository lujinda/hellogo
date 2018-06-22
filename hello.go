package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// 假装自己很忙的一个接口
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		for now+10 > time.Now().Unix() {
		}
		fmt.Fprintln(w, "done")
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		writer := json.NewEncoder(w)
		err := writer.Encode(os.Environ())
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/redis", func(w http.ResponseWriter, r *http.Request) {
		cli, err := redis.DialURL("redis://127.0.0.1:6379/1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if resp, err := redis.String(cli.Do("info", "Server")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, resp)
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Will stop")
		go func() {
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}()
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		log.Println("Received request from", r.RemoteAddr)
		fmt.Fprintf(w, "Received request for URL: %s,\nMy hostname is %q\nV4",
			r.RequestURI, hostname)
	})
	fmt.Println(http.ListenAndServe(":5050", nil))
}
