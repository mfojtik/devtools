package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/data.json") {
		data, err := ioutil.ReadFile("data/data.json")
		if err != nil {
			if _, serveErr := fmt.Fprintf(w, "ERROR: %v", err); serveErr != nil {
				log.Printf("Error serving request: %v", err)
			}
			w.WriteHeader(500)
			return
		}
		_, err = fmt.Fprintln(w, string(data))
		if err != nil {
			if _, serveErr := fmt.Fprintf(w, "ERROR: %v", err); serveErr != nil {
				log.Printf("Error serving request: %v", err)
			}
			w.WriteHeader(500)
			return
		}
	} else {
		serveMux.ServeHTTP(w, r)
	}
}
