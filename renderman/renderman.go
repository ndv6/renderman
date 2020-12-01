package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var skippedResource = []string{
	"googletagmanager",
	"google-analytics",
	"facebook",
	"twitter",
	"windows",
	"service-worker",
}

// ChromeRemote ...
type ChromeRemote struct {
	RemoteURL string
}

// Content ...
type Content struct {
	Title   string
	Content string
}

func (c *Content) Marshal() string {
	content, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(content)
}

func (c *Content) UnMarshal(val []byte) {
	err := json.Unmarshal(val, c)
	if err != nil {
		fmt.Println(err)
	}
}

// FetchHeadless ...
func (c *ChromeRemote) FetchHeadless(url string) (Content, error) {
	var (
		title    string
		content  string
		allocCtx context.Context
		cancel   context.CancelFunc
	)

	if c.RemoteURL != "" {
		allocCtx, cancel = chromedp.NewRemoteAllocator(context.Background(), c.RemoteURL)
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
