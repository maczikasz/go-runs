package main

import (
	"flag"
	"fmt"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/wire/backend"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	config := flag.String("config", "", "Path to the config file (JSON,TOML,YAML,HCL,env, Java properties), if not set config is assumed to be YAML and read from stdin")

	log.Debugf("App started with command line args %s", os.Args)
	flag.Parse()

	if *config == "" {
		viper.SetConfigType("YAML")
		stdInStats, err := os.Stdin.Stat()
		if err != nil {
			panic(err.Error())
		}
		if stdInStats.Size() <= 0 {
			panic("Cannot read from empty StdIn")
		}
		err = viper.ReadConfig(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigFile(*config)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	viper.AutomaticEnv()

	startupContext, cleanup, err := backend.InitializeStartupContext()
	if err != nil {
		panic(err.Error())
	}
	defer cleanup()
	server.StartHttpServer(&wg, startupContext)
	wg.Wait()
}
