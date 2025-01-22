package main

import "net/http"

func main() {
	mx := http.NewServeMux()

	server := http.Server{}

	server.Addr = ":8080"
	server.Handler = mx

	server.ListenAndServe()

}
