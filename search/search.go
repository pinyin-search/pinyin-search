package search

import "pinyin-search/entity"

var MySearch Search

type Search interface {
	Init()
	// tenant 租户 indexName 索引名 doc 文档
	Add(tenant string, indexName string, doc []map[string]interface{}) entity.Result
	Suggestion(tenant string, indexName string, keyword string) entity.Result
}
