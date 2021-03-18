package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/cache"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var hc *Remote
var enableStaticFile = false
var bot = `googlebot|bingbot|adidxBot|yandex|baiduspider|twitterbot|facebookexternalhit|rogerbot|linkedinbot|embedly|quora link preview|showyoubot|outbrain|pinterest\/0\.|pinterestbot|slackbot|slack-imgproxy|vkshare|w3c_validator|whatsapp|collector-agent`
var asset = `\.(js|css|xml|less|png|jpg|jpeg|gif|pdf|doc|txt|ico|rss|zip|mp3|rar|exe|wmv|doc|avi|ppt|mpg|mpeg|tif|wav|mov|psd|ai|xls|mp4|m4a|swf|dat|dmg|iso|flv|m4v|torrent|ttf|woff|svg|eot)`
var appURL string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = godotenv.Load()

	// Echo instance
	e := echo.New()
	hc = &Remote{
		RemoteURL:  os.Getenv("CHROMIUM_URL"),
		BackendURL: os.Getenv("BACKEND_URL"),
	}

	enableStaticFile = os.Getenv("ENABLE_STATIC_SERVE") != ""
	appURL = os.Getenv("APP_URL")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if os.Getenv("CACHE_DRIVER") == "redis" {
		rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			logrus.Error(err)
		}
		cache.Register(cache.DvrRedis, redishashConnect(os.Getenv("REDIS_DSN"), "", rdb))
		cache.DefaultUse(cache.DvrRedis)
	} else {
		cache.Register(cache.DvrFile, cache.NewFileCache("./tmp"))
		cache.DefaultUse(cache.DvrFile)
	}

	e.GET("*", httpHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func httpHandler(ctx echo.Context) error {
	uri := fmt.Sprintf("%s%s", os.Getenv("BACKEND_URL"), ctx.Request().URL)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(uri)))

	if enableStaticFile {
		userAgent := strings.ToLower(ctx.Request().Header.Get("user-agent"))
		matched, _ := regexp.Match(bot, []byte(userAgent))
		if matched {
			fmt.Println("=====fetchAndCacheHeadless====:68", uri)
			return fetchAndCacheHeadless(ctx, uri, hash)
		}

		matched, _ = regexp.Match(asset, []byte(uri))
		if !matched {
			fmt.Println("=====fetchAndCacheHeadless====:74", uri)
			return fetchAndCacheHeadless(ctx, uri, hash)
		}
		fmt.Println("=====proxy====:77", uri)
		hc.Proxy(ctx)
		return nil
	}

	fmt.Println("=====fetchAndCacheHeadless====:82", uri)
	return fetchAndCacheHeadless(ctx, uri, hash)
}

func fetchAndCacheHeadless(ctx echo.Context, uri, hash string) error {
	c := Content{}
	fmt.Println("====user-agent====", ctx.Request().Header.Get("user-agent"))
	key := fmt.Sprintf("%s:%s", appURL, hash)
	fmt.Println("====key", key)
	ctn := cache.Get(key)
	fmt.Println("====ctn == nil", ctn == nil)
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
	fmt.Println("======c=====", c.Header, c.Title, c.Content)
	// replace backend_url to app_url
	c.Content = strings.ReplaceAll(c.Content, os.Getenv("BACKEND_URL"), os.Getenv("APP_URL"))

	// cache

	_ = cache.Set(hash, c.Marshal(), cache.OneDay)
	return ctx.HTML(http.StatusOK, c.Content)
}
