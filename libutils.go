package libutils

import (
	"os"
	"os/exec"
	"os/signal"
	"io"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"sync"
	"strconv"
	"syscall"
	"encoding/json"
)

var (
	Lock = sync.RWMutex{}
	PathFile = os.Args[0]
)

func Atoi(s string) int {
	value, err := strconv.Atoi(s)

	if err != nil {
		panic(err)
	}

	return value
}

func RealPath(name string) string {
	realpath, _ := filepath.Abs(filepath.Dir(PathFile) + "/" + name)

	return realpath
}

func BytesToSize(value float64) string {
	suffixes := []string{
		"B",
		"KB",
		"MB",
		"GB",
	}

	var i int

	for value >= 1024 && i < (len(suffixes) - 1) {
		value = value / 1024
		i++
	}

	return fmt.Sprintf("%.3f %s", value, suffixes[i])
}

func IsCommandExists(file string) bool {
	_, err := exec.LookPath(file)

	return err == nil
}

func CreateFile(name string, s string) error {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

func MakeDir(fullpath string) {
	os.MkdirAll(fullpath, 0700)
}

func CopyFile(source string, destination string, mode int) {
	from, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, os.FileMode(mode))
	if err != nil {
		panic(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		panic(err)
	}
}

func JsonWrite(v interface{}, filename string) {
	bytedata, _ := json.MarshalIndent(v, "", "	")

	MakeDir(filepath.Dir(filename))

	ioutil.WriteFile(filename, bytedata, 0644)
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
