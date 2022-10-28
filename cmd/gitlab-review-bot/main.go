package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/pkg/app"
)

var (
	fConfigPath string
)

func init() {
	flag.StringVar(&fConfigPath, "config", "config/config.yml", "path to yml config file")
}

func main() {
	flag.Parse()

	a, err := app.New(fConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to create app")
		os.Exit(2)
	}

	err = a.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to run app")
		os.Exit(2)
	}
}
