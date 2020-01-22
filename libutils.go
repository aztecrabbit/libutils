package libutils

import (
	"os"
	"os/signal"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"syscall"
	"encoding/json"
)

var (
	PathFile = os.Args[0]
)

func RealPath(name string) string {
	return filepath.Dir(PathFile) + "/" + name
}

func Atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return value
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
	bytedata, _ := json.MarshalIndent(v, "", "	")

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
