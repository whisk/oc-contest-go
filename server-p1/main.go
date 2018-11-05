package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

type person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var cnt int64

func md5hex(text string) string {
	digest := md5.New()
	io.WriteString(digest, text)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	var person person
	err := json.Unmarshal(body, &person)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	response, err := json.Marshal(map[string]string{
		"id":           person.ID,
		"first_name":   person.FirstName + " " + md5hex(person.FirstName),
		"last_name":    person.LastName + " " + md5hex(person.LastName),
		"current_time": time.Now().Format("2006-01-02 15:04:05 -0700"),
		"say":          "go is the best",
	})
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(response)
	if cnt%100 == 0 {
		fmt.Println(cnt)
	}
	cnt++
}

func httpHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		indexHandler(ctx)
	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}

func main() {
	numProcs := flag.Int("n", 1, "Max number of procs")
	port := flag.Int("p", 8080, "Port number")
	host := flag.String("h", "127.0.0.1", "Host")
	flag.Parse()

	runtime.GOMAXPROCS(*numProcs)
	fmt.Println("Started.")
	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), httpHandler)
	if err != nil {
		log.Fatal(err)
	}
}
