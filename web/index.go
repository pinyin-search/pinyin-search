package web

import (
	"net/http"
)

const running = "{\"status\":\"running\"}"

// Index suggest
func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(running))
}
