package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/entity"
	"pinyin-search/search"
)

// Update 更新索引(不存在会新增)
func Update(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")
	dataId := request.Form.Get("dataId")
	data := request.Form.Get("data")

	pyRequest := entity.PinYinRequest{Tenant: tenant, IndexName: indexName, DataId: dataId, Data: data}

	docs, err := pyRequest.GetDocs()
	if err != nil {
		returnJson, _ := json.Marshal(entity.Result{Success: false, Msg: err.Error()})
		writer.Write(returnJson)
		return
	}

	result, err := search.MySearch.Update(tenant+"_"+indexName, dataId, docs)
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)

}
