package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, World!")
		if err != nil {
			return
		}
	})
	err := http.ListenAndServe(":6789", nil)
	if err != nil {
		return
	}
}
