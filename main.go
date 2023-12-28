package main

import (
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

const (
	envHost = "HL_HOST"

	defaultHost = ":8177"
)

var tab = "    "

func handler(w http.ResponseWriter, r *http.Request) {
	color.Cyan("_________________________________________________________________________________________________")
	color.Cyan("| Incoming request")
	color.HiCyan(t(1)+"%s %s", r.Method, r.Host+r.URL.Path)

	color.Yellow(t(1) + "--------------------")
	color.Yellow(t(1) + "Headers")
	for k, v := range r.Header {
		color.HiYellow(t(2)+"%s %v", k, v)
	}

	color.Green(t(1) + "--------------------")
	color.Green(t(1) + "Data")
	data, err := io.ReadAll(r.Body)
	e(err)
	color.HiGreen(t(2)+"%s", data)
	io.WriteString(w, "ACK")
}

func t(num int) string {
	var result = "| "
	for i := 0; i < num; i++ {
		result += tab
	}
	return result
}

func e(err error) {
	if err != nil {
		color.Red(">>> %s", err.Error())
	}
}

func main() {
	http.HandleFunc("/", handler)

	host := os.Getenv(envHost)
	if host == "" {
		host = defaultHost
	}
	color.HiMagenta(">>> Server listening on %s", host)
	panic(http.ListenAndServe(host, nil))
}
