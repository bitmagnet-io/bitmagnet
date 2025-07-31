package env

import (
	"os"
	"os/signal"
)

type osSignalsProvider struct{}

func (osSignalsProvider) Signals(sigs ...os.Signal) <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, sigs...)

	return ch
}

type nopSignalsProvider struct{}

func (nopSignalsProvider) Signals(...os.Signal) <-chan os.Signal {
	return make(chan os.Signal)
}
