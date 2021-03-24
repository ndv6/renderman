package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/go-rod/rod/lib/utils"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("address", ":9222", "the address to listen to")
var quiet = flag.Bool("quiet", false, "silent the log")

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)
}

// a cli tool to launch browser remotely
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	rl := NewRemoteLauncher()
	rl.setup()
	defer rl.cleanup()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		utils.E(err)
	}

	fmt.Println("Remote control url is", "ws://"+l.Addr().String())

	srv := &http.Server{Handler: rl}
	utils.E(srv.Serve(l))
}
