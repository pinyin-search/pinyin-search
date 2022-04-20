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

// add 添加索引
func (meili *MeiliSearch) add(tenant string, indexName string, doc []map[string]interface{}) (entity.Result, error) {
	// An index is where the documents are stored.
	index := meili.Client.Index(tenant + "_" + indexName)

	// 初始化参数
	if _, ok := meili.DistinctMap.LoadOrStore(tenant+"_"+indexName, true); !ok {
		// 主索引为id
		_, err := meili.Client.GetIndex(tenant + "_" + indexName)
		if err, ok := err.(*meilisearch.Error); ok {
			if err.StatusCode == 404 {
				meili.Client.CreateIndex(&meilisearch.IndexConfig{
					Uid:        tenant + "_" + indexName,
					PrimaryKey: "id",
				})
			}
		}

		// guid 可搜索
		index.UpdateFilterableAttributes(&[]string{
			"guid",
		})
		// 结果去重
		index.UpdateDistinctAttribute("value")

	}

	task, err := index.AddDocuments(doc)
	if err != nil {
		log.Printf("添加Index失败。tenant: %s, indexName: %s, err: %s\n", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	task, err = meili.Client.GetTask(task.UID)
	if err != nil {
		log.Printf("添加Index失败。tenant: %s, indexName: %s, err: %s\n", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	return entity.Result{Success: true, Data: doc, Msg: string(task.Status)}, nil
}

// Update 更新索引
func (meili *MeiliSearch) AddUpdate(tenant string, indexName string, guid string, doc []map[string]interface{}) (entity.Result, error) {

	// 通过guid删除旧数据
	meili.Delete(tenant, indexName, guid)

	return meili.add(tenant, indexName, doc)
}

// Delete 删除索引
func (meili *MeiliSearch) Delete(tenant string, indexName string, guid string) (entity.Result, error) {

	index := meili.Client.Index(tenant + "_" + indexName)

	// ResetDistinctAttribute
	t, err := index.ResetDistinctAttribute()
	if err == nil {
		meili.Client.GetTask(t.UID)
	}

	// 通过guid查找旧数据
	resp, err := index.Search("", &meilisearch.SearchRequest{
		Filter: [][]string{
			{"guid = '" + guid + "'"},
		},
		Limit: 1000,
	})

	// 删除旧的索引
	if err == nil {
		deleteIds := make([]string, len(resp.Hits))
		for i, hit := range resp.Hits {
			if hitNew, ok := hit.(map[string]interface{}); ok {
				if id, ok := hitNew["id"].(string); ok {
					deleteIds[i] = id
				}
			}
		}
		index.DeleteDocuments(deleteIds)

		log.Printf("删除 %d 个索引 %s\n", len(deleteIds), resp.Hits)
	} else {
		log.Printf("删除索引失败. Err: %s\n", err)
	}

	return entity.Result{Success: true, Data: nil}, nil
}

// Suggestion 搜索建议
func (meili *MeiliSearch) Suggestion(tenant string, indexName string, keyword string) (entity.Result, error) {
	index := meili.Client.Index(tenant + "_" + indexName)
	task, err := index.UpdateDistinctAttribute("value")
	if err == nil {
		meili.Client.GetTask(task.UID)
	}

	searchRes, err := index.Search(keyword, &meilisearch.SearchRequest{})

	if err != nil {
		log.Printf("搜索失败。tenant: %s, indexName: %s, err: %s", tenant, indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	return entity.Result{Success: true, Data: searchRes.Hits}, nil
}
