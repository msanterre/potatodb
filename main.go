package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	listenAddress = "127.0.0.1:5050"
)

var Data map[string]string

func Error(w http.ResponseWriter, req *http.Request, err string) {
	io.WriteString(w, err)
}

func DBMethod(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.RequestURI, "/db/")

	if key == "" {
		NotFoundError(w)
		return
	}

	switch req.Method {
	case "GET":
		if value, ok := Data[key]; ok {
			io.WriteString(w, value)
		} else {
			NotFoundError(w)
		}
	case "POST", "PUT", "PATCH":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ServerError(w, err)
			log.Println(err)
			return
		}
		Data[key] = string(body)

		w.WriteHeader(200)
	case "DELETE":
		delete(Data, key)
		w.WriteHeader(200)
	}
}

func NotFoundError(w http.ResponseWriter) {
	w.WriteHeader(404)
	io.WriteString(w, "Not found")
}

func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	io.WriteString(w, err.Error())
}

func main() {
	Data = make(map[string]string)

	http.HandleFunc("/db/", DBMethod)

	fmt.Printf("Listening on %s\n", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
