package search

import (
	"fmt"
	"log"
	"os"
	"pinyin-search/entity"
	"strings"
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
func (meili *MeiliSearch) Add(indexName string, docs []entity.Doc) (entity.Result, error) {
	// An index is where the documents are stored.
	index := meili.Client.Index(indexName)

	// 初始化参数
	if _, ok := meili.DistinctMap.LoadOrStore(indexName, true); !ok {
		// 主索引为id
		_, err := meili.Client.GetIndex(indexName)
		if err, ok := err.(*meilisearch.Error); ok {
			if err.StatusCode == 404 {
				meili.Client.CreateIndex(&meilisearch.IndexConfig{
					Uid:        indexName,
					PrimaryKey: "id",
				})
			}
		}

		// 结果去重
		index.UpdateDistinctAttribute("value")

		// 设置只能key能搜索
		sa := make([]string, 2)
		sa = append(sa, "key")
		sa = append(sa, "value")
		index.UpdateSearchableAttributes(&sa)
	}

	taskDoc, err := index.AddDocuments(docs)
	if err != nil {
		log.Printf("添加Index失败。indexName: %s, err: %s\n", indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	taskIdx, err := meili.Client.GetTask(taskDoc.TaskUID)
	if err != nil {
		log.Printf("添加Index失败。indexName: %s, err: %s\n", indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}
	log.Printf("添加Index成功。indexName: %s\n", indexName)

	return entity.Result{Success: true, Data: docs, Msg: string(taskIdx.Status)}, nil
}

// Update 更新索引
func (meili *MeiliSearch) Update(indexName string, dataId string, docs []entity.Doc) (entity.Result, error) {

	// 通过guid删除旧数据
	meili.Delete(indexName, []string{dataId})

	return meili.Add(indexName, docs)
}

// Delete 删除索引
func (meili *MeiliSearch) Delete(indexName string, dataIds []string) (entity.Result, error) {

	_, err := meili.Client.GetIndex(indexName)
	if err != nil {
		return entity.Result{Success: false, Msg: fmt.Sprintf("删除索引失败! Err: %s", err.Error())}, err
	}

	index := meili.Client.Index(indexName)

	var deleteAllIds []string

	// 通过dataId查找旧数据
	for _, dataId := range dataIds {
		resp, err := index.Search(dataId, &meilisearch.SearchRequest{
			AttributesToRetrieve: []string{"id"},
			Limit:                1,
		})
		if err == nil {
			deleteIds := make([]string, resp.EstimatedTotalHits)
			var i int64
			for i = 0; i < resp.EstimatedTotalHits; i++ {
				deleteIds[i] = fmt.Sprintf("%s_%d", dataId, i)
			}
			deleteAllIds = append(deleteAllIds, deleteIds...)
		}
	}
	index.DeleteDocuments(deleteAllIds)
	return entity.Result{Success: true, Msg: fmt.Sprintf("删除 %d 条索引", len(deleteAllIds))}, nil
}

// DeleteAll 删除indexName下的全部索引
func (meili *MeiliSearch) DeleteAll(indexName string) (entity.Result, error) {

	task, err := meili.Client.Index(indexName).DeleteAllDocuments()

	// 删除旧的索引
	if err == nil {
		log.Printf("删除 %s 的所有Index成功\n", indexName)
		return entity.Result{Success: true, Msg: string(task.Status)}, nil
	} else if err, ok := err.(*meilisearch.Error); ok && err.StatusCode != 404 {
		log.Printf("删除 %s 的所有Index失败, Err: %s\n", indexName, err)
		return entity.Result{Success: false, Msg: err.Error(), Data: nil}, err
	}
	return entity.Result{Success: true}, nil
}

// Suggest 搜索建议
func (meili *MeiliSearch) Suggest(indexName string, keyword string) (entity.Result, error) {
	if strings.TrimSpace(keyword) == "" {
		return entity.Result{Success: true, Data: make([]string, 0)}, nil
	}

	searchRes, err := meili.Client.Index(indexName).Search(keyword, &meilisearch.SearchRequest{})

	if err != nil {
		log.Printf("搜索失败。indexName: %s, Err: %s", indexName, err)
		return entity.Result{Success: false, Msg: err.Error()}, err
	}

	return entity.Result{Success: true, Data: searchRes.Hits}, nil
}
