package app

import (
	config "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	HumanReadableLog bool   `config:"human_readable_log"`
	GitlabToken      string `config:"gitlab_token"`
	SlackToken       string `config:"slack_token"`

	Mongo struct {
		Host string `config:"host"`
		Port int    `config:"port"`
		User string `config:"user"`
		Pass string `config:"pass"`
		DB   string `config:"db"`
	} `config:"mongo"`
}

func (a *App) initConfig(configPath string) error {
	_ = godotenv.Load()

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yamlv3.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
	})

	err := config.LoadFiles(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	err = config.Decode(&a.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to decode config")
	}

	return nil
}
