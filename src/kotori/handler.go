package main

import (
	"net/http"
	"fmt"
)

func Pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Pong!")
}
