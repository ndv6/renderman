package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var skippedResource = []string{
	"googletagmanager",
	"google-analytics",
	"facebook",
	"twitter",
	"windows",
	"service-worker",
}

var skippedResourceType = map[network.ResourceType]interface{}{
	network.ResourceTypeFont:  nil,
	network.ResourceTypeImage: nil,
}

// Remote ...
type Remote struct {
	RemoteURL  string
	BackendURL string
}

// Content ...
type Content struct {
	Title   string
	Content string
	Header  http.Header
}

// Marshal ...
func (c *Content) Marshal() string {
	content, err := json.Marshal(c)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return string(content)
}

// UnMarshal ...
func (c *Content) UnMarshal(val []byte) {
	err := json.Unmarshal(val, c)
	if err != nil {
		logrus.Error(err)
	}
}

// Proxy is for serve as proxy
func (rt *Remote) Proxy(c echo.Context) {
	uri, err := url.Parse(rt.BackendURL)
	if err != nil {
		logrus.Error(err)
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

// FetchHeadless ...
func (rt *Remote) FetchHeadless(url string) (Content, error) {
	var (
		title    string
		content  string
		allocCtx context.Context
		cancel   context.CancelFunc
	)

	if rt.RemoteURL != "" {
		allocCtx, cancel = chromedp.NewRemoteAllocator(context.Background(), rt.RemoteURL)
		defer cancel()
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("disable-extensions", false),
			chromedp.Flag("disable-web-security", true),
			chromedp.Flag("auto-open-devtools-for-tabs", true),
		)

		allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
	}

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// intercept
	chromedp.ListenTarget(ctx, interceptAndWaitUntilLoaded(ctx, cancel))

	err := chromedp.Run(ctx,
		fetch.Enable(),
		network.Enable(),
		page.Enable(),
		page.SetLifecycleEventsEnabled(true),
		chromedp.Navigate(url),
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &content),
	)
	if err != nil {
		return Content{}, err
	}

	return Content{Title: title, Content: content}, nil
}

func interceptAndWaitUntilLoaded(ctx context.Context, cancelFn context.CancelFunc) func(ev interface{}) {
	var frameID string
	var init bool
	return func(event interface{}) {
		switch ev := event.(type) {
		case *fetch.EventRequestPaused:
			c := chromedp.FromContext(ctx)
			ctx := cdp.WithExecutor(ctx, c.Target)
			for _, skip := range skippedResource {
				if strings.Contains(ev.Request.URL, skip) {
					go fetch.FailRequest(ev.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
					return
				}
			}
			if _, ok := skippedResourceType[ev.ResourceType]; ok {
				go fetch.FailRequest(ev.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
				return
			}
			go fetch.ContinueRequest(ev.RequestID).Do(ctx)
		case *page.EventFrameStartedLoading:
			frameID = ev.FrameID.String()
		case *page.EventLifecycleEvent:
			switch ev.Name {
			case "init":
				init = ev.FrameID.String() == frameID
			case "networkIdle":
				if init && ev.FrameID.String() == frameID {
					cancelFn()
				}
			}
		}
	}
}
