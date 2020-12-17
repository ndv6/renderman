package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

const flagKeepUserDataDir = "rod-keep-user-data-dir"

// RemoteLauncher ...
type RemoteLauncher struct {
	Logger utils.Logger
}

// NewRemoteLauncher instance
func NewRemoteLauncher() *RemoteLauncher {
	return &RemoteLauncher{
		Logger: utils.LoggerQuiet,
	}
}

func (p *RemoteLauncher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		p.launch(w, r)
		return
	}

	p.defaults(w, r)
}

func newLauncher() *launcher.Launcher {
	l := launcher.New()
	l.Set("disable-web-security")
	l.Set("no-sandbox")
	l.Set("disable-setuid-sandbox")
	l.Set("disable-dev-shm-usage")
	l.Set("disable-accelerated-2d-canvas")
	l.Set("disable-gpu")
	l.Set("disable-notifications")
	l.Set("rod-keep-user-data-dir")

	return l
}

func (p *RemoteLauncher) defaults(w http.ResponseWriter, _ *http.Request) {
	l := newLauncher()
	utils.E(w.Write(l.JSON()))
}

func (p *RemoteLauncher) launch(w http.ResponseWriter, r *http.Request) {
	l := newLauncher()

	options := r.Header.Get(launcher.HeaderName)
	if options != "" {
		l.Flags = nil
		utils.E(json.Unmarshal([]byte(options), l))
	}

	u := l.Leakless(false).MustLaunch()
	defer func() {
		l.Kill()
		p.Logger.Println("Killed PID:", l.PID())

		if _, has := l.Get(flagKeepUserDataDir); !has {
			l.Cleanup()
			dir, _ := l.Get("user-data-dir")
			p.Logger.Println("Removed", dir)
		}
	}()

	parsedURL, err := url.Parse(u)
	utils.E(err)

	p.Logger.Println("Launch", u, l.FormatArgs())
	defer p.Logger.Println("Close", u)

	parsedWS, err := url.Parse(u)
	utils.E(err)
	parsedURL.Path = parsedWS.Path

	httputil.NewSingleHostReverseProxy(toHTTP(*parsedURL)).ServeHTTP(w, r)
}

func toHTTP(u url.URL) *url.URL {
	newURL := u
	if newURL.Scheme == "ws" {
		newURL.Scheme = "http"
	} else if newURL.Scheme == "wss" {
		newURL.Scheme = "https"
	}
	return &newURL
}
