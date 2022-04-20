package search

import (
	"log"
	"os"
	"pinyin-search/entity"
	"sync"

	"github.com/meilisearch/meilisearch-go"
)

const MEILISEARCH_HOST_ENV = "MEILISEARCH_HOST"
const MEILISEARCH_APIKEY_ENV = "MEILISEARCH_APIKEY"

type MeiliSearch struct {
	Client      *meilisearch.Client
	DistinctMap sync.Map
}

func (meili *MeiliSearch) Init() {
	meili.Client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   os.Getenv(MEILISEARCH_HOST_ENV),
		APIKey: os.Getenv(MEILISEARCH_APIKEY_ENV),
	})
}

// Add 添加索引
func (meili *MeiliSearch) Add(tenant string, indexName string, doc []map[string]interface{}) (entity.Result, error) {
	// An index is where the documents are stored.
	index := meili.Client.Index(tenant + "_" + indexName)

	// 结果去重
	if _, ok := meili.DistinctMap.LoadOrStore(tenant+"_"+indexName, true); !ok {
		index.UpdateDistinctAttribute("value")
	}

	task, err := index.AddDocuments(doc)
	if err != nil {
		log.Printf("添加Index失败。tenant: %s, indexName: %s, err: %s\n", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	return entity.Result{Success: true, Data: doc, Msg: string(task.Status)}, nil
}

// Suggestion 搜索建议
func (meili *MeiliSearch) Suggestion(tenant string, indexName string, keyword string) (entity.Result, error) {
	searchRes, err := meili.Client.Index(tenant+"_"+indexName).Search(keyword, &meilisearch.SearchRequest{})
	if err != nil {
		log.Printf("搜索失败。tenant: %s, indexName: %s, err: %s", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	return entity.Result{Success: true, Data: searchRes.Hits}, nil
}
