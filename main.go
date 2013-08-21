package main

import (
	goflags "github.com/jessevdk/go-flags"
	pm "github.com/politician/go-proxymate/proxymate"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	ProxyAddr  string `short:"c" long:"client" value-name:"HOST" description:"The address to which the client will connect." optional:"true"`
	ListenAddr string `short:"s" long:"server" value-name:"HOST" description:"The address on which the server will listen." required:"true"`
	Rate       int    `short:"r" long:"rate" description:"The number of operations per second initiated by the client."`
}

func main() {
	// Parse the command line.
	opts := parseArgs()

	// Start up the server, so that the client has something to talk to.
	server, err := pm.ListenAndServe(opts.ListenAddr, &pm.SampleHandler{})
	if err != nil {
		log.Fatal(err)
	}

	// Dial the server, and start pinging it.
	client, err := pm.DialAndPing(opts.ProxyAddr, opts.Rate)
	if err != nil {
		log.Fatal(err)
	}

	AwaitSignals(func() {
		client.Stop()
		server.Stop()
	})

	log.Println("goodbye")
}

func AwaitSignals(shutdown func()) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Println("CTRL-C to exit...")

	for {
		// Block until we receive a signal.
		sig := <-ch
		log.Println("Got: ", sig.String())

		switch sig {

		// TODO - Handle other signals that don't just stop the
		// process immediately.

		// SIGQUIT should exit gracefully.
		case syscall.SIGQUIT:
			shutdown()
			return

		// SIGTERM should exit.
		case syscall.SIGTERM, syscall.SIGINT:
			shutdown()
			return
		}
	}
}

// Parses the command-line arguments, and validates them.
func parseArgs() *Options {
	opts := &Options{}

	_, err := goflags.Parse(opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.ProxyAddr == "" {
		opts.ProxyAddr = opts.ListenAddr
	}

	switch {
	case opts.Rate < 1:
		opts.Rate = 2
	case opts.Rate > 10:
		opts.Rate = 10
	}

	return opts
}
