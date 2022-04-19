package search

import (
	"log"
	"os"
	"pinyin-search/entity"

	"github.com/meilisearch/meilisearch-go"
)

const MEILISEARCH_HOST_ENV = "MEILISEARCH_HOST"
const MEILISEARCH_APIKEY_ENV = "MEILISEARCH_APIKEY"

type MeiliSearch struct {
	Client *meilisearch.Client
}

func (meili *MeiliSearch) Init() {
	meili.Client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   os.Getenv(MEILISEARCH_HOST_ENV),
		APIKey: os.Getenv(MEILISEARCH_APIKEY_ENV),
	})
}

func (meili *MeiliSearch) Add(tenant string, indexName string, doc []map[string]interface{}) entity.Result {
	// An index is where the documents are stored.
	index := meili.Client.Index(tenant + "_" + indexName)
	task, err := index.AddDocuments(doc)
	if err != nil {
		log.Printf("添加Index失败。tenant: %s, indexName: %s, err: %s\n", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}
	}

	t, err := index.GetTask(task.UID)
	if err != nil {
		log.Printf("添加Index失败(获取task失败)。tenant: %s, indexName: %s, err: %s", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}
	}
	success := t.Status == meilisearch.TaskStatusEnqueued || t.Status == meilisearch.TaskStatusSucceeded || t.Status == meilisearch.TaskStatusProcessing
	return entity.Result{Success: success, Data: doc}
}

func (meili *MeiliSearch) Suggestion(tenant string, indexName string, keyword string) entity.Result {
	searchRes, err := meili.Client.Index(tenant+"_"+indexName).Search(keyword, &meilisearch.SearchRequest{})
	if err != nil {
		log.Printf("搜索失败。tenant: %s, indexName: %s, err: %s", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}
	}

	return entity.Result{Success: true, Data: searchRes.Hits}
}
