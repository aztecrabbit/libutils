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
	"bufio"
	"strconv"
	"syscall"
	"strings"
	"runtime"
	"encoding/json"
)

var (
	Lock = sync.RWMutex{}
	Stdin = bufio.NewReader(os.Stdin)
	PathFile, _ = os.Executable()
)

func Atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return value
}

func PaddingLeft(value string, s string, count int) string {
	s = strings.Repeat(s, count)
	if len(value) >= len(s) {
		return value
	}

	return s[:len(s) - len(value)] + value
}

func Input(s string) string {
	fmt.Printf(s)
	value, _ := Stdin.ReadString('\n')

	return strings.TrimSuffix(value, "\n")
}

func RealPath(name string) string {
	return filepath.Dir(PathFile) + "/" + name
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

func CopyFile(source string, destination string) {
	from, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
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

func KillProcess(p *os.Process) {
	switch runtime.GOOS {
		case "windows":
			p.Kill()
		default:
			p.Signal(syscall.SIGTERM)
	}
}

//

type InterruptHandler struct {
	Done chan bool
	Handle func()
}

func (i *InterruptHandler) Start() {
	i.Done = make(chan bool)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func () {
		<- ch
		if i.Handle != nil {
			i.Handle()
		}
		i.Done <- true
	}()
}

func (i *InterruptHandler) Wait() {
	<- i.Done
	os.Exit(0)
}
