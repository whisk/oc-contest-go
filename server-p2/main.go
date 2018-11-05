package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/valyala/fasthttp"
)

type person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type response struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CurrentTime string `json:"current_time"`
	Say         string
}

var cnt int64
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func md5hex(text string) string {
	digest := md5.New()
	io.WriteString(digest, text)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	var person person
	err := person.UnmarshalJSON(body)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	z, err := response{
		ID:          person.ID,
		FirstName:   person.FirstName + " " + md5hex(person.FirstName),
		LastName:    person.LastName + " " + md5hex(person.LastName),
		CurrentTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
		Say:         "go is the best",
	}.MarshalJSON()

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(z)
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
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Writing cpu profile to %s...\n", *cpuprofile)
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			fmt.Printf("Cpu profile completed.\n")
		}()
	}

	runtime.GOMAXPROCS(*numProcs)
	fmt.Println("Started.")
	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), httpHandler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Printf("Wrote mem profile to %s.\n", *cpuprofile)
	}
}
