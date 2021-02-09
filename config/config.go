package config

import "encoding/json"

type Config struct {
	Port        int    `envconfig:"PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
}

func (c Config) Fields() map[string]interface{} {
	b, _ := json.Marshal(c)

	var fields map[string]interface{}
	json.Unmarshal(b, &fields)
	return fields
}
