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
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const say = "go is the best"
const timeFormat = "2006-01-02 15:04:05 -0700"
const contentType = "application/json"

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
var serverDateStr = ""

func md5hex(text string) []byte {
	digest := md5.New()
	io.WriteString(digest, text)
	return digest.Sum(nil)
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	var person person
	err := person.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.Error("Unmarshal error: "+err.Error(), fasthttp.StatusBadRequest)
		return
	}

	var first strings.Builder
	first.WriteString(person.FirstName)
	first.WriteString(" ")
	first.Write(md5hex(person.FirstName))

	var last strings.Builder
	first.WriteString(person.LastName)
	first.WriteString(" ")
	first.Write(md5hex(person.LastName))

	res, err := response{
		ID:          person.ID,
		FirstName:   first.String(),
		LastName:    last.String(),
		CurrentTime: serverDateStr,
		Say:         say,
	}.MarshalJSON()

	if err != nil {
		ctx.Error("Marshal error: "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(res)
	if cnt%1000 == 0 {
		fmt.Println(cnt)
	}
	cnt++
}

func httpHandler(ctx *fasthttp.RequestCtx) {
	indexHandler(ctx)
}

func main() {
	numProcs := flag.Int("n", 1, "Max number of procs")
	port := flag.Int("p", 8080, "Port number")
	host := flag.String("h", "127.0.0.1", "Host")
	concurrency := flag.Int("c", 256*1024, "fasthttp concurrency")
	keepalive := flag.Bool("k", true, "Keepalive")
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

	go func() {
		for {
			serverDateStr = time.Now().Format(timeFormat)
			time.Sleep(time.Second)
		}
	}()

	runtime.GOMAXPROCS(*numProcs)
	go func() {
		s := fasthttp.Server{
			DisableHeaderNamesNormalizing: true,
			NoDefaultServerHeader:         true,
			Logger:                        nil,
			Handler:                       httpHandler,
			DisableKeepalive:              !*keepalive,
			Concurrency:                   *concurrency,
		}
		fmt.Printf("Started %+v.\n", s)
		err := s.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port))
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
