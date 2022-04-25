package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/entity"
	"pinyin-search/search"
)

// DeleteAll 删除Index下所有的索引
func DeleteAll(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")

	if indexName == "" {
		writer.WriteHeader(400)
		returnJson, _ := json.Marshal(entity.Result{Success: false, Msg: "请求参数不完整"})
		writer.Write(returnJson)
		return
	}

	result, err := search.MySearch.DeleteAll(tenant+"_"+indexName)
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)

}
