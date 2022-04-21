package search

import "pinyin-search/entity"

var MySearch Search

type Search interface {
	Init()
	// Add/Update 通过dataId更新索引
	AddUpdate(indexName string, dataId string, doc []map[string]interface{}) (entity.Result, error)
	Delete(indexName string, dataId string) (entity.Result, error)
	Suggestion(indexName string, keyword string) (entity.Result, error)
}
