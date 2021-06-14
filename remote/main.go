package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("address", ":9222", "the address to listen to")
var quiet = flag.Bool("quiet", false, "silence the log")
var allowAllPath = flag.Bool("allow-all", false, "allow all path set by the client")

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

var l *launcher.Launcher

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	m := launcher.NewManager()

	if l == nil {
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

		// l.Headless(false)
		l.Leakless(true)
		l.UserDataDir("/tmp")
		// l.MustLaunch()
	}

	m.Defaults = func(_ http.ResponseWriter, _ *http.Request) *launcher.Launcher {
		return l
	}

	if !*quiet {
		m.Logger = logrus.New()
	}

	if *allowAllPath {
		m.BeforeLaunch = func(l *launcher.Launcher, rw http.ResponseWriter, r *http.Request) {}
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		utils.E(err)
	}

	if !*quiet {
		fmt.Println("rod-manager listening on:", l.Addr().String())
	}

	srv := &http.Server{Handler: m}
	utils.E(srv.Serve(l))
}
