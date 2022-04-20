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

	result, err := search.MySearch.Suggestion(tenant, indexName, data)
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)
}
