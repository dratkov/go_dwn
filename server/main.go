package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", GoServer)
	http.ListenAndServe(":8080", nil)
}

func GoServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=")
	fmt.Fprintf(w, "- Go Go Go -")
}
