package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/search"
)

// Delete 删除索引
func Delete(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")
	dataId := request.Form.Get("dataId")

	result, err := search.MySearch.Delete(tenant+"_"+indexName, dataId)
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)

}
