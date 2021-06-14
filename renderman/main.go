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
	"time"

	"github.com/gookit/cache"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var (
	renderman        *Renderman
	enableStaticFile = false
	bot              = `googlebot|bingbot|adidxBot|yandex|baiduspider|twitterbot|facebookexternalhit|rogerbot|linkedinbot|embedly|quora link preview|showyoubot|outbrain|pinterest\/0\.|pinterestbot|slackbot|slack-imgproxy|vkshare|w3c_validator|whatsapp|collector-agent`
	asset            = `\.(js|css|xml|less|png|jpg|jpeg|gif|pdf|doc|txt|ico|rss|zip|mp3|rar|exe|wmv|doc|avi|ppt|mpg|mpeg|tif|wav|mov|psd|ai|xls|mp4|m4a|swf|dat|dmg|iso|flv|m4v|torrent|ttf|woff|svg|eot)`
	hashName         string
	compressr        compressor
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logrus.Info(runtime.GOOS)

	_ = godotenv.Load()

	// Echo instance
	e := echo.New()
	renderman = newRenderman(os.Getenv("BACKEND_URL"), os.Getenv("CHROMIUM_URL"))
	compressr = newZstdCompressor()

	e.HideBanner = true

	enableStaticFile = os.Getenv("ENABLE_STATIC_SERVE") != ""
	hashName = strings.ReplaceAll(os.Getenv("APP_URL"), ":", "")

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
	if !enableStaticFile {
		return fetchAndCacheHeadless(ctx, uri)
	}

	userAgent := strings.ToLower(ctx.Request().Header.Get("user-agent"))
	matched, _ := regexp.Match(bot, []byte(userAgent))
	if matched {
		return fetchAndCacheHeadless(ctx, uri)
	}

	matched, _ = regexp.Match(asset, []byte(uri))
	if !matched {
		return fetchAndCacheHeadless(ctx, uri)
	}
	renderman.Proxy(ctx)
	return nil
}

func fetchFromRemote(uri string) (c Content, err error) {
	c, err = renderman.FetchHeadless(uri)
	if err != nil {
		logrus.Error(err)
		return
	}
	// replace backend_url to app_url
	c.Content = strings.ReplaceAll(c.Content, os.Getenv("BACKEND_URL"), os.Getenv("APP_URL"))
	return
}

func storeToCache(key string, c Content, dur time.Duration) {
	byts := []byte(c.Marshal())
	byts = compressr.Compress(byts)
	_ = cache.Set(key, byts, dur)
}

func fetchAndCacheHeadless(ctx echo.Context, uri string) error {
	var (
		hash = fmt.Sprintf("%x", md5.Sum([]byte(uri)))
		key  = fmt.Sprintf("%s:%s", hashName, hash)
		ctn  interface{}
		byts []byte
		c    Content
	)
	if ctx.Request().Header.Get("user-agent") == "collector-agent" || !cache.Has(key) {
		c, err := fetchFromRemote(uri)
		if err != nil {
			logrus.Error(err)
			return ctx.HTML(http.StatusInternalServerError, "")
		}
		storeToCache(key, c, cache.OneDay)
		return ctx.HTML(http.StatusOK, c.Content)
	}
	ctn = cache.Get(key)
	byts = ctn.([]byte)
	byts, err := compressr.Decompress(byts)
	if err != nil {
		logrus.Error(err)
		return ctx.HTML(http.StatusInternalServerError, "")
	}
	c.UnMarshal(byts)
	return ctx.HTML(http.StatusOK, c.Content)
}
