package web

import (
	"encoding/json"
	"net/http"
	"pinyin-search/entity"
	"pinyin-search/search"
)

type Batch []entity.PinYinRequest

// UpdateBatch 批量更新索引
func UpdateBatch(writer http.ResponseWriter, request *http.Request) {
	batch := &Batch{}
	err := json.NewDecoder(request.Body).Decode(batch)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	batchMap := make(map[string][]entity.Doc)
	for _, pyRequest := range *batch {
		docs, err := pyRequest.GetDocs()
		if err == nil {
			name := pyRequest.Tenant + "_" + pyRequest.IndexName
			batchMap[name] = append(batchMap[name], docs...)
			// delete first
			search.MySearch.Delete(name, pyRequest.DataId)
		}
	}

	successNum := 0
	for k, v := range batchMap {
		_, err := search.MySearch.Add(k, v)
		if err == nil {
			successNum = successNum + 1
		}
	}

	if successNum == len(batchMap) {
		returnJson, _ := json.Marshal(entity.Result{Success: true, Msg: "批量新增成功"})
		writer.Write(returnJson)
	} else {
		returnJson, _ := json.Marshal(entity.Result{Success: true, Msg: "批量新增部分成功, err:" + err.Error()})
		writer.Write(returnJson)
	}

}
