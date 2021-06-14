package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var skippedResource = []string{
	"googletagmanager",
	"google-analytics",
	"facebook",
	"twitter",
	"service-worker",
}

var skippedResourceType = map[proto.NetworkResourceType]interface{}{
	proto.NetworkResourceTypeFont:  nil,
	proto.NetworkResourceTypeImage: nil,
	proto.NetworkResourceTypeMedia: nil,
}

var xhrAllowRequests = map[proto.NetworkResourceType]interface{}{
	proto.NetworkResourceTypePreflight: nil,
	proto.NetworkResourceTypeXHR:       nil,
	proto.NetworkResourceTypeFetch:     nil,
}

// Remote ...
type Renderman struct {
	browser    *rod.Browser
	remoteUrl  string
	backendURL string
}

func newRenderman(backendUrl, remoteUrl string) *Renderman {
	renderman := &Renderman{backendURL: backendUrl, remoteUrl: remoteUrl}
	renderman.browser = renderman.GetBrowser()

	return renderman
}

func (renderman *Renderman) GetBrowser() *rod.Browser {
	if renderman.browser == nil {
		var l *launcher.Launcher
		if renderman.remoteUrl == "" {
			path, _ := launcher.LookPath()
			logrus.Info("chrome path :", path)

			l = launcher.New().Bin(path)

			l.Set("disable-web-security")
			l.Set("no-sandbox")
			l.Set("disable-setuid-sandbox")
			l.Set("disable-dev-shm-usage")
			l.Set("disable-accelerated-2d-canvas")
			l.Set("disable-gpu")
			l.Set("disable-notifications")
			l.Set("rod-keep-user-data-dir")
			// renderman.browser = rod.New().ControlURL(l.MustLaunch()).MustConnect()
		} else {
			// u := launcher.MustResolveURL(renderman.remoteUrl)
			l = launcher.MustNewManaged(renderman.remoteUrl)
		}
		renderman.browser = rod.New().ControlURL(l.MustLaunch()).MustConnect()
	}

	return renderman.browser
}

// Proxy is for serve as proxy
func (renderman *Renderman) Proxy(c echo.Context) {
	uri, err := url.Parse(renderman.backendURL)
	if err != nil {
		logrus.Error("unable parse url ", err)
	}

	config := middleware.ProxyConfig{
		Skipper:    middleware.DefaultSkipper,
		ContextKey: "target",
	}

	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := uri.String()
		c.Set("_error", echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err)))
		logrus.Error(err)
	}
	proxy.Transport = config.Transport
	proxy.ModifyResponse = config.ModifyResponse

	req := c.Request()
	res := c.Response()

	req.URL.Host = uri.Host
	req.URL.Scheme = uri.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = uri.Host

	if req.Header.Get(echo.HeaderXRealIP) == "" || c.Echo().IPExtractor != nil {
		req.Header.Set(echo.HeaderXRealIP, c.RealIP())
	}
	if req.Header.Get(echo.HeaderXForwardedProto) == "" {
		req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
	}

	proxy.ServeHTTP(res, req)
}

func (renderman *Renderman) FetchHeadless(url string) (Content, error) {
	browser := renderman.GetBrowser()
	page, err := browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		renderman.browser = nil
		return Content{}, err
	}

	router := page.HijackRequests()
	router.MustAdd("*", func(ctx *rod.Hijack) {
		// skip unused fetch request
		for _, skip := range skippedResource {
			if strings.Contains(fmt.Sprint(ctx.Request.URL()), skip) {
				logrus.WithFields(logrus.Fields{"type": "resource", "policies": "blocked"}).Info(ctx.Request.URL())
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				return
			}
		}

		// Ensure xhr request is loaded
		if _, ok := xhrAllowRequests[ctx.Request.Type()]; ok {
			log := logrus.WithFields(logrus.Fields{"type": strings.ToLower(fmt.Sprint(ctx.Request.Type())), "policies": "allow"})
			log.Info(ctx.Request.Method(), " ", ctx.Request.URL())
			ctx.MustLoadResponse()
			// if err := ctx.LoadResponse(http.DefaultClient, true); err != nil {
			// 	log.Info("complete")
			// } else {
			// 	log.Errorf("failed %s", err.Error())
			// }
			// ctx.ContinueRequest(&proto.FetchContinueRequest{})
		}

		// skip asset resouce image and font
		if _, ok := skippedResourceType[ctx.Request.Type()]; ok {
			logrus.WithFields(logrus.Fields{"type": strings.ToLower(fmt.Sprint(ctx.Request.Type())), "policies": "blocked"}).Info(ctx.Request.URL())
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	// wait := page.MustWaitNavigation() // wait until network is almost idle
	// wait := page.MustWaitRequestIdle()
	wait := page.WaitRequestIdle(1*time.Second, nil, []string{})
	if err := page.Navigate(url); err != nil {
		_ = renderman.browser.Close()
		renderman.browser = nil

		logrus.Error("unable to navigate page ", err)
		return Content{}, err
	}
	wait()

	html, err := page.MustElement("html").HTML()
	if err != nil {
		_ = renderman.browser.Close()
		renderman.browser = nil

		logrus.Error("unable to get page ", err)
		return Content{}, err
	}

	if err := page.Close(); err != nil {
		logrus.Error("unable to close page ", err)
	}

	return Content{Content: html}, nil
}
