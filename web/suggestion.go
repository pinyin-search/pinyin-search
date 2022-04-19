package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/search"
)

// Suggestion Suggestion
func Suggestion(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")
	data := request.Form.Get("data")
	returnJson, _ := json.Marshal(search.MySearch.Suggestion(tenant, indexName, data))
	writer.Write(returnJson)
}
