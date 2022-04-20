package search

import "pinyin-search/entity"

var MySearch Search

type Search interface {
	Init()
	// Add/Update 通过guid更新索引
	AddUpdate(tenant string, indexName string, guid string, doc []map[string]interface{}) (entity.Result, error)
	Delete(tenant string, indexName string, guid string) (entity.Result, error)
	Suggestion(tenant string, indexName string, keyword string) (entity.Result, error)
}
