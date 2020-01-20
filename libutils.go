package libutils

import (
	"os"
	"os/signal"
	"syscall"
)

type InterruptHandler struct {
	Handle func()
}

func (i *InterruptHandler) Start() {
    channel := make(chan os.Signal, 2)
    signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

    go func () {
    	<- channel
    	if i.Handle != nil {
    		i.Handle()
	    }
    	os.Exit(0)
    }()
}
