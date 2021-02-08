package logger

import "sync"

type Logger struct {
	config
	itemChan chan Item
	done     sync.WaitGroup
}

func New(configurers ...Configurer) *Logger {
	cfg := defaultConfig()
	for _, configurer := range configurers {
		configurer(&cfg)
	}

	logger := &Logger{
		config:   cfg,
		itemChan: make(chan Item, cfg.BufferSize),
	}

	logger.done.Add(1)
	go logger.worker()
	return logger
}

func (l *Logger) Close() {
	close(l.itemChan)
	l.done.Wait()
}
