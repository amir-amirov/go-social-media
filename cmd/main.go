package main

import "net/http"

type server struct {
	addr string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func main() {

	s := &server{
		addr: ":8080",
	}
	http.ListenAndServe((*s).addr, s)
}
