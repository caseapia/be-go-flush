package config

import (
	"log"
	"os"
)

type Config struct {
	DiscordClientID              string
	DiscordClientSecret          string
	DiscordLocalRedirectURI      string
	DiscordProductionRedirectURI string
	AppMode                      string
}

func Load() *Config {
	cfg := &Config{
		DiscordClientID:              os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret:          os.Getenv("DISCORD_CLIENT_SECRET"),
		DiscordLocalRedirectURI:      os.Getenv("DISCORD_LOCAL_REDIRECT_URI"),
		DiscordProductionRedirectURI: os.Getenv("DISCORD_PRODUCTION_REDIRECT_URI"),
		AppMode:                      os.Getenv("APP_MODE"),
	}

	validate(cfg)

	return cfg
}

func validate(cfg *Config) {
	if cfg.DiscordClientID == "" {
		log.Fatal("DISCORD_CLIENT_ID is not set")
	}
	if cfg.DiscordClientSecret == "" {
		log.Fatal("DISCORD_CLIENT_SECRET is not set")
	}
	if cfg.DiscordLocalRedirectURI == "" {
		log.Fatal("DISCORD_LOCAL_REDIRECT_URI is not set")
	}
	if cfg.DiscordProductionRedirectURI == "" {
		log.Fatal("DISCORD_PRODUCTION_REDIRECT_URI is not set")
	}
	if cfg.AppMode == "" {
		log.Fatal("APP_MODE is not set")
	}
}

func (c *Config) DiscordRedirectURI() string {
	if c.AppMode == "production" {
		return c.DiscordProductionRedirectURI
	}

	return c.DiscordLocalRedirectURI
}
