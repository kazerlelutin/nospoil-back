package handlers

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "hello")
}
