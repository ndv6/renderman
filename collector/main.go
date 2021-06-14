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

var log *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	log = logrus.WithFields(logrus.Fields{
		"app":       "renderman",
		"component": "collector",
	})
}

func main() {
	time.Sleep(1 * time.Minute)
	runtime.GOMAXPROCS(runtime.NumCPU())
	godotenv.Load()
	uri := os.Getenv("APP_URL")

	if os.Getenv("COLLECTOR_ENABLE") != "" {
		log.Info("collector initialized")

		cron := cron.New()
		scheduler := os.Getenv("COLLECTOR_SCHEDULER")
		if scheduler == "" {
			scheduler = "0 */6 * * *"
		}

		cron.AddFunc(scheduler, func() {
			collect(uri)
		})

		collect(uri)
		go func() {
			log.Info("start cron scheduler")
			cron.Start()
		}()

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig

	} else {
		logrus.Info("collector is disabled")
	}
}

func collect(uri string) {
	log.Info("cron start at ", time.Now())
	c := colly.NewCollector(colly.UserAgent("collector-agent"))
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
	})

	log.Infof("started collect")
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		time.Sleep(30 * time.Millisecond)
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info("collector:", r.URL)
	})

	c.Visit(uri)
	c.Wait()
	log.Info("completed collect")
}
