package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// delay collector before main app start
	time.Sleep(2 * time.Minute)

	runtime.GOMAXPROCS(runtime.NumCPU())
	godotenv.Load()
	uri := os.Getenv("APP_URL")

	if os.Getenv("COLLECTOR_ENABLE") == "true" {
		cron := cron.New()

		logrus.Info("collector agent is enabled")

		scheduler := os.Getenv("COLLECTOR_SCHEDULER")
		if scheduler == "" {
			scheduler = "0 */6 * * *"
		}

		cron.AddFunc(scheduler, func() {
			logrus.Info("cron start at ", time.Now())
			collect(uri)
		})

		collect(uri)
		go func() {
			logrus.Info("start cron scheduler")
			cron.Start()
		}()

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig

	} else {
		logrus.Info("collector agent is disabled")
	}
}

func collect(uri string) {
	c := colly.NewCollector(colly.UserAgent("collector-agent"))
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		logrus.Info("collector:", r.URL)
	})

	c.Visit(uri)
}
