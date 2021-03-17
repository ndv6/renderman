package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/go-rod/rod/lib/utils"
)

var addr = flag.String("address", ":9222", "the address to listen to")
var quiet = flag.Bool("quiet", false, "silent the log")

// a cli tool to launch browser remotely
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	rl := NewRemoteLauncher()
	rl.setup()
	defer rl.cleanup()
	if !*quiet {
		rl.Logger = log.New(os.Stdout, "", 0)
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		utils.E(err)
	}

	fmt.Println("Remote control url is", "ws://"+l.Addr().String())

	srv := &http.Server{Handler: rl}
	utils.E(srv.Serve(l))
}
