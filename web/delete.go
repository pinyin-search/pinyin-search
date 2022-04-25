package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/entity"
	"pinyin-search/search"
)

// Delete 删除索引
func Delete(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")
	dataId := request.Form.Get("dataId")

	if dataId == "" || indexName == "" {
		writer.WriteHeader(400)
		returnJson, _ := json.Marshal(entity.Result{Success: false, Msg: "请求参数不完整"})
		writer.Write(returnJson)
		return
	}

	result, err := search.MySearch.Delete(tenant+"_"+indexName, []string{dataId})
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)

}
