package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Fieldser interface {
	Fields() map[string]interface{}
}

type Fields map[string]interface{}

func (f Fields) Fields() map[string]interface{} {
	return f
}

func Field(key string, value interface{}) Fields {
	return Fields{key: value}
}

type Item struct {
	Level     Level
	Message   string
	Fields    []Fieldser
	Timestamp time.Time
}

func (l *Logger) Log(level Level, message string, fields ...Fieldser) {
	now := time.Now().UTC()
	go func() { l.itemChan <- Item{Level: level, Message: message, Fields: fields, Timestamp: now} }()
}

func (l *Logger) Err(err error, message string, fields ...Fieldser) {
	l.Error(message, append(fields, Field("error", err.Error()))...)
}

func (l *Logger) worker() {
	defer l.done.Done()

	for item := range l.itemChan {
		if item.Level < l.Level {
			continue
		}

		payload := make(map[string]interface{})
		for _, fields := range item.Fields {
			for key, value := range fields.Fields() {
				if key != "" {
					payload[key] = value
				}
			}
		}

		payload["message"] = item.Message
		payload["level"] = strings.ToLower(item.Level.String())
		payload["time"] = item.Timestamp.Format(l.TimeFormat)
		if l.HumanFormat {
			if _, err := fmt.Fprint(l.Writer, humanFormat(payload)); err != nil {
				fmt.Printf("failed to human encode payload: payload: %v, error: %s\n", payload, err)
			}
		} else {
			if err := json.NewEncoder(l.Writer).Encode(payload); err != nil {
				fmt.Printf("failed to encode payload: payload: %v, error: %s\n", payload, err)
				return
			}
		}
	}
}

func humanFormat(payload map[string]interface{}) string {
	var others []string
	for key, value := range payload {
		if key == "time" || key == "level" || key == "message" {
			continue
		}

		valueString := fmt.Sprint(value)
		if strings.Contains(valueString, " ") {
			valueString = fmt.Sprintf("(%s)", valueString)
		}

		others = append(others, fmt.Sprintf("%s=%v", key, valueString))
	}

	otherString := strings.Join(others, " ")
	if otherString != "" {
		otherString = fmt.Sprintf("; %s", otherString)
	}

	level := strings.ToUpper(payload["level"].(string))
	return fmt.Sprintf("%v [%v] %v%s\n", payload["time"], level, payload["message"], otherString)
}
