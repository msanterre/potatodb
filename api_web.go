package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func DBMethod(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.RequestURI, "/db/")

	if key == "" {
		WebNotFoundError(w)
		return
	}

	switch req.Method {
	case "GET":
		if value, ok := DB.Data[key]; ok {
			io.WriteString(w, value)
		} else {
			WebNotFoundError(w)
		}
	case "POST", "PUT", "PATCH":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			WebServerError(w, err)
			log.Println(err)
			return
		}
		DB.Data[key] = string(body)

		w.WriteHeader(200)
	case "DELETE":
		delete(DB.Data, key)
		w.WriteHeader(200)
	}
}

func WebNotFoundError(w http.ResponseWriter) {
	w.WriteHeader(404)
	io.WriteString(w, "Not found")
}

func WebServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	io.WriteString(w, err.Error())
}

func WebStart(listenAddress string) {
	http.HandleFunc("/db/", DBMethod)

	log.Printf("Web listening on %s\n", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
