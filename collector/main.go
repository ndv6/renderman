package main

import (
	"os"
	"runtime"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	godotenv.Load()
	uri := os.Getenv("APP_URL")

	if os.Getenv("COLLECTOR_ENABLE") != "" {
		logrus.Info("collector initialized")
		c := colly.NewCollector(colly.UserAgent("collector-agent"))
		c.Limit(&colly.LimitRule{
			Parallelism: 1,
		})

		logrus.Infof("started collect")
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			time.Sleep(30 * time.Millisecond)
			e.Request.Visit(e.Attr("href"))
		})
		c.OnRequest(func(r *colly.Request) {
			logrus.Info("collector:", r.URL)
		})

		c.Visit(uri)
		c.Wait()
		logrus.Info("completed collect")
	} else {
		logrus.Info("collector is disabled")
	}
}
