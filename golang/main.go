package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/sync/syncmap"
)

// Server as http server
type Server struct {
	// m  map[string]string
	m syncmap.Map
}

func (s *Server) start() {
	fmt.Println("Starting server on 8081 ...")

	http.HandleFunc("/", s.handler)
	http.ListenAndServe(":8081", nil)

	fmt.Println("Done")
}

func (s *Server) handler(w http.ResponseWriter, req *http.Request) {

	paths := strings.SplitN(req.URL.Path, "/", 3)
	if len(paths) < 2 {
		http.Error(w, "Must provide a key in the path", http.StatusBadRequest)
		return
	}

	key := paths[1]
	switch req.Method {
	case "GET":
		v, _ := s.m.Load(key)
		vStr, _ := v.(string)
		fmt.Fprintf(w, vStr)
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Cannot read request body", http.StatusBadRequest)
		}
		value := string(body)
		s.m.Store(key, value)
		fmt.Fprintf(w, value)
		// fmt.Println(value)
	default:
		http.Error(w, "Only GET and POST methods are supported.", http.StatusMethodNotAllowed)
	}
}

func main() {
	server := Server{m: syncmap.Map{}}
	server.start()
}
