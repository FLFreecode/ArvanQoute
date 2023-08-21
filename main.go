package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	apiServer "github.com/arvan/qoute/api"
	"github.com/arvan/qoute/config"
)

var (
	configFile     string
	cert           string
	key            string
	addr           string
	verbosityLevel int
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -cert {cert file}")
	fmt.Println("      -key {key file}")
	fmt.Println("      -a {listen addr}")
	fmt.Println("      -h (show help info)")
	fmt.Println("      -v {0-10} (verbosity level, default 0)")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	flag.StringVar(&configFile, "c", "config.yml", "config file")

	if !config.Load(configFile) {
		os.Exit(1)
	}

	if verbosityLevel < 0 {
		verbosityLevel = config.Get().Log.Level
	}
	zerolog.SetGlobalLevel(zerolog.Level(verbosityLevel))
	serv := apiServer.NewAppServer(config.Get())

	log.Print("Arvan Qoute service is ON")
	if serv.TracingCloser != nil {
		defer serv.TracingCloser.Close()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)
	select {
	case err := <-serv.ListenAndServe():

		panic(err)
	case <-sigCh:
		fmt.Println("Shutdown service...")
		os.Exit(1)
	}
}
