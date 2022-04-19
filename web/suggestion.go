package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/search"
)

// Suggestion Suggestion
func Suggestion(writer http.ResponseWriter, request *http.Request) {
	indexName := request.FormValue("indexName")
	data := request.FormValue("data")
	returnJson, _ := json.Marshal(search.MySearch.Suggestion("test", indexName, data))
	writer.Write(returnJson)
}
