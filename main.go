package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	envHost = "HL_HOST"

	defaultHost = ":8177"
)

var (
	host     = defaultHost
	output   = "terminal"
	filepath = ""

	outputs = map[string]func(r *http.Request, data []byte){
		"terminal":      outputTerminal,
		"terminal-json": outputTerminalJSON,
		"file":          outputFile,
		"file-json":     outputFileJSON,
	}
)

func init() {
	flag.StringVar(&host, "host", defaultHost, "Defines the host of the service. (default: :8177)")
	flag.StringVar(&output, "output", "terminal", "Defines the output(s), separated with semicolon. (default: terminal)")
	flag.StringVar(&filepath, "filepath", "", "Defines the filepath to the file output.")
	flag.Parse()
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	e(err)
	for _, o := range strings.Split(output, ";") {
		outputs[o](r, data)
	}
	_, err = io.WriteString(w, "ACK")
	e(err)
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

func outputTerminal(r *http.Request, data []byte) {
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
	color.HiGreen(t(2)+"%s", data)
}

func outputTerminalJSON(r *http.Request, data []byte) {
	var j = requestToJSON(r, data)
	result, err := json.Marshal(j)
	e(err)
	fmt.Println(string(result))
}

func outputFile(r *http.Request, data []byte) {
	f := openFile()
	if f == nil {
		return
	}
	var builder strings.Builder
	builder.WriteString("_________________________________________________________________________________________________\n")
	builder.WriteString("| Incoming request\n")
	builder.WriteString(fmt.Sprintf(t(1)+"%s %s\n", r.Method, r.Host+r.URL.Path))

	builder.WriteString(t(1) + "--------------------\n")
	builder.WriteString(t(1) + "Headers\n")
	for k, v := range r.Header {
		builder.WriteString(fmt.Sprintf(t(2)+"%s %v\n", k, v))
	}

	builder.WriteString(t(1) + "--------------------\n")
	builder.WriteString(t(1) + "Data\n")
	builder.WriteString(fmt.Sprintf(t(2)+"%s\n", data))
	_, err := fmt.Fprint(f, builder.String())
	e(err)
}

func outputFileJSON(r *http.Request, data []byte) {
	f := openFile()
	if f == nil {
		return
	}
	var j = requestToJSON(r, data)
	result, err := json.Marshal(j)
	e(err)
	_, err = fmt.Fprintln(f, string(result))
	e(err)
}

func openFile() *os.File {
	var f *os.File
	var err error
	if _, err = os.Stat(filepath); err == nil {
		f, err = os.Open(filepath)
		if err != nil {
			e(err)
			return nil
		}

	} else if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(filepath)
		if err != nil {
			e(err)
			return nil
		}
	} else {
		e(err)
		return nil
	}
	return f
}

func requestToJSON(r *http.Request, data []byte) map[string]interface{} {
	var result = make(map[string]interface{})
	result["method"] = r.Method
	result["path"] = r.Host + r.URL.Path
	result["headers"] = r.Header
	var dataContainer = make(map[string]interface{})
	err := json.Unmarshal(data, &dataContainer)
	e(err)
	result["body"] = dataContainer

	return result
}

func t(num int) string {
	var tab = "    "
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
