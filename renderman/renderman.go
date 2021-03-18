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

func (c *Content) Compress() []byte {
	str := c.Marshal()

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
	chromedp.ListenTarget(ctx, func(event interface{}) {
		go func() {
			switch ev := event.(type) {
			case *fetch.EventRequestPaused:
				c := chromedp.FromContext(ctx)
				ctx := cdp.WithExecutor(ctx, c.Target)

				for _, skip := range skippedResource {
					if strings.Contains(ev.Request.URL, skip) {
						fetch.FailRequest(ev.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
						break
					}
				}

				if ev.ResourceType == network.ResourceTypeImage {
					fetch.FailRequest(ev.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
				} else {
					fetch.ContinueRequest(ev.RequestID).Do(ctx)
				}
			}
		}()
	})

	err := chromedp.Run(ctx,
		fetch.Enable(),
		network.Enable(),
		enableLifeCycleEvents(),
		navigate(url),
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &content),
	)
	if err != nil {
		return Content{}, err
	}

	return Content{Title: title, Content: content}, nil
}

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}
		err = page.SetLifecycleEventsEnabled(true).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

// Navigate is an action that navigates the current frame.
func navigate(urlstr string) chromedp.NavigateAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		_, _, _, err := page.Navigate(urlstr).Do(ctx)
		if err != nil {
			return err
		}
		return waitLoaded(ctx)
	})
}

// waitLoaded blocks until a target receives a Page.loadEventFired.
func waitLoaded(ctx context.Context) error {
	ch := make(chan struct{})
	lctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == "networkAlmostIdle" {
				cancel()
				close(ch)
			}
		}
	})

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
