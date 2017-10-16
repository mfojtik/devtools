package main

import (
	"bufio"
	"log"
	"os"

	"github.com/mfojtik/devtools/logs/lde/glog"
	"github.com/mfojtik/devtools/logs/lde/request"
)

func main() {
	reader, err := os.Open("/Users/mfojtik/Downloads/origin-master-api.service")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	scanner := bufio.NewScanner(reader)
	l := &glog.Line{}
	for scanner.Scan() {
		ok, err := l.Extract(scanner.Bytes())
		if !ok {
			continue
		}
		if err != nil {
			log.Printf("Failed to extract data: %v", err)
			continue
		}
		r := &request.Line{}
		r.Extract(l.Message)
		if len(r.Verb) > 0 {
			log.Printf("v: %q", string(r.Verb))
			log.Printf("p: %q", string(r.Path))
			log.Printf("d: %q", string(r.Duration))
			log.Printf("c: %d", r.StatusCode)
			log.Printf("---------")
		}
	}
}
