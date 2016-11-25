package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var ExitChan chan os.Signal

func init() {
	ExitChan = make(chan os.Signal)
	signal.Notify(ExitChan, syscall.SIGINT, syscall.SIGTERM)
}
