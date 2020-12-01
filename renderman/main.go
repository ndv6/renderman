package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/gookit/cache"
	"github.com/gookit/cache/redis"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Echo instance
	e := echo.New()
	hc := ChromeRemote{RemoteURL: os.Getenv("CHROMIUM_URL")}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if os.Getenv("CACHE_DRIVER") == "redis" {
		rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			fmt.Println(err)
		}
		cache.Register(cache.DvrRedis, redis.Connect(os.Getenv("REDIS_DSN"), "", rdb))
		cache.DefaultUse(cache.DvrRedis)
	} else {
		cache.Register(cache.DvrFile, cache.NewFileCache("./tmp"))
		cache.DefaultUse(cache.DvrFile)
	}

	e.GET("*", func(ctx echo.Context) error {
		uri := fmt.Sprintf("%s%s", os.Getenv("BACKEND_URL"), ctx.Request().URL)
		hash := fmt.Sprintf("%x", md5.Sum([]byte(uri)))

		//
		ctn := cache.Get(hash)
		if ctn != nil {
			c := Content{}
			c.UnMarshal([]byte(fmt.Sprintf("%s", ctn)))

			if ctx.Request().Header.Get("user-agent") != "collector-agent" {
				return ctx.HTML(http.StatusOK, c.Content)
			}
		}

		c, err := hc.FetchHeadless(uri)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.HTML(http.StatusInternalServerError, "")
		}

		// cache
		cache.Set(hash, c.Marshal(), cache.OneDay)

		return ctx.HTML(http.StatusOK, c.Content)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
