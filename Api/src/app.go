package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("In a hello handler")
	fmt.Fprintf(w, "Hello", "from handler")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":7000", nil)
}
