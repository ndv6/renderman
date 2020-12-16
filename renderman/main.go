package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"

	"github.com/gookit/cache"
	"github.com/gookit/cache/redis"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var hc *Remote
var enableStaticFile = false

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	godotenv.Load()

	// Echo instance
	e := echo.New()
	hc = &Remote{
		RemoteURL:  os.Getenv("CHROMIUM_URL"),
		BackendURL: os.Getenv("BACKEND_URL"),
	}

	if os.Getenv("ENABLE_STATIC_SERVE") != "" {
		enableStaticFile = true
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if os.Getenv("CACHE_DRIVER") == "redis" {
		rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			logrus.Error(err)
		}
		cache.Register(cache.DvrRedis, redis.Connect(os.Getenv("REDIS_DSN"), "", rdb))
		cache.DefaultUse(cache.DvrRedis)
	} else {
		cache.Register(cache.DvrFile, cache.NewFileCache("./tmp"))
		cache.DefaultUse(cache.DvrFile)
	}

	e.GET("*", httpHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func httpHandler(ctx echo.Context) error {
	renderman := false
	uri := fmt.Sprintf("%s%s", os.Getenv("BACKEND_URL"), ctx.Request().URL)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(uri)))

	bot := `googlebot|bingbot|yandex|baiduspider|twitterbot|facebookexternalhit|rogerbot|linkedinbot|embedly|quora link preview|showyoubot|outbrain|pinterest\/0\.|pinterestbot|slackbot|Slack-ImgProxy|Slackbot-LinkExpanding|vkShare|W3C_Validator|whatsapp|collector-agent`
	matched, _ := regexp.Match(bot, []byte(ctx.Request().Header.Get("user-agent")))
	if matched == true {
		renderman = true
	}

	asset := `\.(js|css|xml|less|png|jpg|jpeg|gif|pdf|doc|txt|ico|rss|zip|mp3|rar|exe|wmv|doc|avi|ppt|mpg|mpeg|tif|wav|mov|psd|ai|xls|mp4|m4a|swf|dat|dmg|iso|flv|m4v|torrent|ttf|woff|svg|eot)`
	matched, _ = regexp.Match(asset, []byte(uri))
	if matched == true {
		renderman = false
	}

	if enableStaticFile == true {
		if renderman == true {
			return fetchAndCacheHeadless(ctx, uri, hash)
		}
		hc.Proxy(ctx)
		return nil
	}

	return fetchAndCacheHeadless(ctx, uri, hash)
}

func fetchAndCacheHeadless(ctx echo.Context, uri, hash string) error {
	c := Content{}
	ctn := cache.Get(hash)
	if ctn != nil {
		c.UnMarshal([]byte(fmt.Sprintf("%s", ctn)))

		// if user agent is not collector-agent then we can return from cache
		if ctx.Request().Header.Get("user-agent") != "collector-agent" {
			return ctx.HTML(http.StatusOK, c.Content)
		}
	}

	c, err := hc.FetchHeadless(uri)
	if err != nil {
		logrus.Error(err)
		return ctx.HTML(http.StatusInternalServerError, "")
	}

	// cache
	cache.Set(hash, c.Marshal(), cache.OneDay)
	return ctx.HTML(http.StatusOK, c.Content)
}
