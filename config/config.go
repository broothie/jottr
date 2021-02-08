package config

type Config struct {
	Port        int    `default:"8080"`
	Environment string `default:"development"`
	SaveDelayMs int    `default:"750"`
}
