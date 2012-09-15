package gpm_index_server

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/all", handler_all)
	http.HandleFunc("/publish", handler_publish)
}

func handler_all(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func handler_publish(w http.ResponseWriter, r *http.Request) {

}
