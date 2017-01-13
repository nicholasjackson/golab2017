package main

import "net/http"

func main() {
	http.DefaultServeMux.HandleFunc("/hello", handle)
	http.ListenAndServe(":9090", http.DefaultServeMux)
}

func handle(rw http.ResponseWriter, r *http.Request) {

}
