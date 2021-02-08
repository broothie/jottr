package logger

import (
	"io"
	"os"
	"time"
)

type config struct {
	Level       Level
	Writer      io.Writer
	BufferSize  uint
	TimeFormat  string
	HumanFormat bool
}

func defaultConfig() config {
	return config{
		Level:      Debug,
		Writer:     os.Stdout,
		BufferSize: 64,
		TimeFormat: time.RFC3339Nano,
	}
}

type Configurer func(*config)

func UseLevel(level Level) Configurer {
	return func(c *config) { c.Level = level }
}

func UseWriter(writer io.Writer) Configurer {
	return func(c *config) { c.Writer = writer }
}

func UseBufferSize(size uint) Configurer {
	return func(c *config) { c.BufferSize = size }
}

func UseTimeFormat(timeFormat string) Configurer {
	return func(c *config) { c.TimeFormat = timeFormat }
}

func UseHumanFormatSetting(setting bool) Configurer {
	return func(c *config) { c.HumanFormat = setting }
}

func UseHumanFormat() Configurer {
	return UseHumanFormatSetting(true)
}
