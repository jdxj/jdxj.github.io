package main

import "net/http"

func main() {
	r, _ := http.NewRequest()
	r.Close
}
