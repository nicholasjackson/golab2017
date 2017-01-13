package main

import "net/http"

func main() {
	http.DefaultServeMux.HandleFunc("/currency", handle)
	http.ListenAndServe(":9091", http.DefaultServeMux)
}

func handle(rw http.ResponseWriter, r *http.Request) {

}
