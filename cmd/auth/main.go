package main

import (
	"flag"
	"fmt"
	"github.com/maczikasz/go-runs/internal/auth"
	authconfig "github.com/maczikasz/go-runs/internal/wire/auth"
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

	flag.Parse()

	if *config == "" {
		viper.SetConfigType("YAML")
		err := viper.ReadConfig(os.Stdin)
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

	context := authconfig.InitializeAuthContext()
	auth.StartHttpServer(&wg, context)
	wg.Wait()

}
