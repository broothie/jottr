package config

type Config struct {
	Port        int    `envconfig:"PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	SaveDelayMs int    `envconfig:"SAVE_DELAY_MS" default:"750"`
}
