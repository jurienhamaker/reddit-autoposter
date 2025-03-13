package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"jurien.dev/reddit-recurring/discord"
	"jurien.dev/reddit-recurring/reddit"

	externalReddit "github.com/vartanbeno/go-reddit/v2/reddit"
)

type (
	Config struct {
		Version           float64
		DiscordWebhookUrl string `toml:"discord_webhook_url"`
		Posts             map[string]reddit.PostConfig
	}
)

var config Config

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		slog.Info("Loading .env file")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if os.Getenv("ENV") == "production" {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	configFile := "config.toml"
	if _, err := os.Stat(configFile); err != nil {
		slog.Error("Couldn't find config.toml", "err", err)
		panic(err)
	}

	meta, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		slog.Error("Couldn't parse config.toml", "err", err)
		panic(err)
	}

	keys := meta.Undecoded()
	for _, key := range keys {
		slog.Warn("Unknown configuration key", "key", key)
	}

	slog.Info(fmt.Sprintf("Loaded configuration %.1f", config.Version))
	if config.DiscordWebhookUrl != "" {
		slog.Info("Enabled discord webhook")
	}

	if err := run(); err != nil {
		slog.Error("Something wen't wrong setting up cron jobs", "err", err)
		panic(err)
	}
}

func run() (err error) {
	c := cron.New()
	for name, post := range config.Posts {
		subReddit := strings.TrimLeft(post.Reddit, "/r/")
		slog.Info(fmt.Sprintf("Setting cron for post %s", name), "cron", post.Cron, "reddit", subReddit)

		_, err = c.AddFunc(post.Cron, createCRONFunc(name, post))
		if err != nil {
			return
		}
	}

	c.Run()
	return
}

func createCRONFunc(name string, post reddit.PostConfig) func() {
	return func() {
		posted, err := reddit.Post(name, post)
		if err != nil {
			slog.Error("Something wen't wrong posting", "err", err)
		}

		if config.DiscordWebhookUrl != "" {
			var fullPost *externalReddit.PostAndComments

			if posted != nil {
				fullPost, _ = reddit.GetPost(posted.ID)
			}

			err = discord.Post(config.DiscordWebhookUrl, name, post, fullPost, err)
			if err != nil {
				slog.Error("Something wen't wrong sending webhook", "err", err)
			}
		}
	}
}
