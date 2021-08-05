package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"strings"
	"time"
)

type Args struct {
	target      string
	httpMethod  string
	header      []string
	concurrency int
}

type Response struct {
	resp *http.Response
	took float64
}

func main() {

	args := parseArgs()

	fmt.Printf("%+v\n", args)

	request(args)

}
func request(args Args) {

	run := 0

	client := &http.Client{}
	req, _ := http.NewRequest(args.httpMethod, args.target, nil)
	for _, h := range args.header {
		parts := strings.Split(h, ":")
		if len(parts) == 2 {
			req.Header.Set(strings.Trim(parts[0], " "), strings.Trim(parts[1], " "))
		}
	}

	log := make(map[int64][]Response)
	var stamp time.Time
	sem := make(chan int, args.concurrency)
	for {

		clock := time.Now()
		if stamp.Unix() < clock.Unix() {

			fmt.Printf("[%s] %s %s\n", stamp.Format(time.RFC3339), args.httpMethod, args.target)
			for _, msg := range log[stamp.Unix()] {
				if msg.resp.StatusCode >= 200 && msg.resp.StatusCode < 300 {
					color.Set(color.FgGreen)
				} else {
					color.Set(color.FgRed)
				}
				fmt.Printf("\t%s %d - took %f\n", msg.resp.Proto, msg.resp.StatusCode, msg.took)
				color.Unset()
			}
			delete(log, stamp.Unix())
			stamp = clock
		}

		sem <- 1
		go func() {

			response, _ := client.Do(req)
			elapsed := time.Since(clock)
			run++

			log[stamp.Unix()] = append(log[stamp.Unix()], Response{
				resp: response,
				took: elapsed.Seconds(),
			})
			<-sem
		}()
	}
}

func parseArgs() Args {
	var target = flag.String("target", "http://127.0.0.1", "Target Url")
	var httpMethod = flag.String("method", "GET", "Http Method")
	var header = flag.String("header", "", "http header")
	var concurency = flag.Int("c", 1, "concurrency")
	var help = flag.Bool("help", false, "print this help")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var normalizedHeader []string
	for _, s := range strings.Split(*header, ",") {
		normalizedHeader = append(normalizedHeader, strings.Trim(s, " "))
	}

	return Args{
		target:      *target,
		httpMethod:  strings.ToUpper(*httpMethod),
		header:      normalizedHeader,
		concurrency: *concurency,
	}
}
