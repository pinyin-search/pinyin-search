package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pinyin-search/entity"
	"pinyin-search/search"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-ego/gse"

	"github.com/mozillazg/go-pinyin"
)

var pyArgs = pinyin.NewArgs()
var pyArgsFirst = pinyin.NewArgs()
var seg gse.Segmenter

var regWord, _ = regexp.Compile("[a-zA-Z]")

func init() {
	pyArgsFirst.Style = pinyin.FirstLetter
	// 加载默认字典
	seg.LoadDict()
}

// AddUpdate 新增或更新索引
func AddUpdate(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	tenant := request.Form.Get("tenant")
	indexName := request.Form.Get("indexName")
	dataId := request.Form.Get("dataId")
	data := request.Form.Get("data")

	if  data == "" || dataId == "" || indexName == "" {
		returnJson, _ := json.Marshal(entity.Result{Success: false, Msg: "错误的数据"})
		writer.WriteHeader(400)
		writer.Write(returnJson)
		return
	}

	words := seg.Cut(data, true)
	indexes := make(map[string]string)

	// 全拼字母
	indexes[strings.Join(pinyin.LazyPinyin(data, pyArgs), "")] = data

	// 全拼第一个字母
	firstLetterAll := strings.Join(pinyin.LazyPinyin(data, pyArgsFirst), "")
	if firstLetterAll != "" {
		indexes[firstLetterAll] = data
	}

	// 分词
	for _, word := range words {
		if utf8.RuneCountInString(word) > 1 {
			// 首字母
			firstLetter := strings.Join(pinyin.LazyPinyin(word, pyArgsFirst), "")
			if firstLetter != "" {
				indexes[firstLetter] = word
			}

			// 拼音
			pyAll := ""
			for _, py := range pinyin.LazyPinyin(word, pyArgs) {
				if py != "" {
					indexes[py] = word
					pyAll = pyAll + py
				}
			}
			if pyAll != "" {
				indexes[pyAll] = word
			}
		}
	}

	// 字母单词
	wordsEnglish := ""
	for _, word := range strings.Fields(data) {
		if word != "" {
			wordsEnglish = wordsEnglish + " " + strings.Join(regWord.FindAllString(word, -1), "")
		}
	}
	wordsEnglish = strings.TrimSpace(wordsEnglish)
	if wordsEnglish != "" {
		indexes[wordsEnglish] = wordsEnglish
	}

	var doc []map[string]interface{}
	idx := 0
	for k, v := range indexes {
		// id 为 dataId_index
		doc = append(doc, map[string]interface{}{"id": fmt.Sprintf("%s_%d", dataId, idx), "key": k, "value": v})
		idx = idx + 1
	}

	result, err := search.MySearch.AddUpdate(tenant+"_"+indexName, dataId, doc)
	if err != nil {
		writer.WriteHeader(400)
	}
	returnJson, _ := json.Marshal(result)
	writer.Write(returnJson)

}
