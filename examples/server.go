package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	port = flag.Int("port", 8080, "assign the port of server listen")
)

func main() {
	flag.Parse()
	http.HandleFunc("/foo", normal)
	http.HandleFunc("/issue26", issue26_fileupload)

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func normal(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseForm()
	timeout, err := strconv.Atoi(req.FormValue("timeout"))
	if err == nil && timeout > 0 {
		fmt.Println("timeout=", timeout)
		time.Sleep(time.Duration(timeout) * time.Second)
	}

	ip := req.RemoteAddr
	w.Header().Add("X-Test", "true")
	_, _ = fmt.Fprintf(w, "bar: %d, %s", 200, ip)
}

// issue26_fileupload running examples/fasthttp-reverse-proxy/proxy.go as proxy
// curl -X POST "http://localhost:8081/local/issue26" -F "file=@/path/to/fasthttp-reverse-proxy/LICENSE" -v
//
func issue26_fileupload(w http.ResponseWriter, req *http.Request) {
	const maxFileSize = 1 * 1024 * 1024 // 1MB
	// pull in the uploaded file into memory

	fmt.Printf("%+v", req.Header)

	_ = req.ParseMultipartForm(maxFileSize)
	fd, fh, err := req.FormFile("file")
	checkError(err)
	defer fd.Close()

	fmt.Println(fh.Filename, " was received")

	_, _ = fmt.Fprintf(w, fh.Filename)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
