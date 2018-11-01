package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

type person struct {
	ID        string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var cnt int64

func md5hex(text string) string {
	digest := md5.New()
	io.WriteString(digest, text)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	var person person
	err = json.Unmarshal(body, &person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	response, err := json.Marshal(map[string]string{
		"id":           person.ID,
		"first_name":   person.FirstName + " " + md5hex(person.FirstName),
		"last_name":    person.LastName + " " + md5hex(person.LastName),
		"current_time": time.Now().UTC().String(),
		"say":          "go is the best",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	if cnt%100 == 0 {
		fmt.Println(cnt)
	}
	cnt++
}

func main() {
	runtime.GOMAXPROCS(2)
	http.HandleFunc("/", indexHandler)
	fmt.Println("Started")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
