package search

import "pinyin-search/entity"

var MySearch Search

type Search interface {
	Init()
	Add(indexName string, docs []entity.Doc) (entity.Result, error)
	// Update 通过dataId更新索引
	Update(indexName string, dataId string, docs []entity.Doc) (entity.Result, error)
	Delete(indexName string, dataId[] string) (entity.Result, error)
	DeleteAll(indexName string) (entity.Result, error)
	Suggest(indexName string, keyword string) (entity.Result, error)
}
