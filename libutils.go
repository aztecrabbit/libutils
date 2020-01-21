package libutils

import (
	"os"
	"os/signal"
	"io/ioutil"
	"path/filepath"
	"syscall"
	"encoding/json"
)

var (
	PathFile string
)

func RealPath(name string) string {
	return filepath.Dir(PathFile) + "/" + name
}

func JsonReadWrite(filename string, v interface{}, vd interface{}) {
	r, err := os.Open(filename)
	if err != nil {
		JsonWrite(vd, filename)
		r, _ = os.Open(filename)
	}

	bytedata, _ := ioutil.ReadAll(r)

	json.Unmarshal(bytedata, v)
}

func JsonWrite(v interface{}, filename string) {
	bytedata, _ := json.MarshalIndent(v, "", "    ")

	ioutil.WriteFile(filename, bytedata, 0644)
}

//

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
    	os.Exit(1)
    }()
}
