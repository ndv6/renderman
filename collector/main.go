package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// delay collector before main app start
	time.Sleep(5 * time.Second)

	runtime.GOMAXPROCS(runtime.NumCPU())
	godotenv.Load()
	uri := os.Getenv("APP_URL")

	if os.Getenv("COLLECTOR_ENABLE") == "true" {
		cron := cron.New()

		fmt.Println("collector agent is enabled")
		// c.Visit(uri)

		cron.AddFunc("* * * * *", func() {
			fmt.Print("cront start ", time.Now())
			collect(uri)
		})

		collect(uri)
		go func() {
			fmt.Println("start cron")
			cron.Start()
		}()

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig

	} else {
		fmt.Println("collector agent is disabled")
	}
}

func collect(uri string) {
	c := colly.NewCollector(colly.UserAgent("collector-agent"))
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("collector:", r.URL)
	})

	c.Visit(uri)
}
