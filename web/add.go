package web

import (
	"encoding/json"
	"hash/crc32"
	"net/http"
	"pinyin-search/search"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mozillazg/go-pinyin"
	"github.com/yanyiwu/gojieba"
)

var pyArgs = pinyin.NewArgs()
var pyArgsFirst = pinyin.NewArgs()
var jieba = gojieba.NewJieba()

var regWord, _ = regexp.Compile("[a-zA-Z]")

func init() {
	pyArgsFirst.Style = pinyin.FirstLetter
}

// Add add
func Add(writer http.ResponseWriter, request *http.Request) {
	indexName := request.FormValue("indexName")
	data := request.FormValue("data")
	words := jieba.Cut(data, true)
	indexes := make(map[string]string)

	// 全部
	indexes[strings.Join(pinyin.LazyPinyin(data, pyArgs), "")] = data

	// 分词
	for _, word := range words {
		if utf8.RuneCountInString(word) > 1 {
			// 首字母
			firstLetter := strings.Join(pinyin.LazyPinyin(word, pyArgsFirst), "")
			if firstLetter != "" {
				indexes[firstLetter] = word
			}
			// 拼音
			for _, w := range pinyin.LazyPinyin(word, pyArgs) {
				if w != "" {
					indexes[w] = word
				}
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
	indexes[wordsEnglish] = wordsEnglish

	var doc []map[string]interface{}
	for k, v := range indexes {
		doc = append(doc, map[string]interface{}{"id": int(crc32.ChecksumIEEE([]byte(k + v))), "key": k, "value": v})
	}

	returnJson, _ := json.Marshal(search.MySearch.Add("test", indexName, doc))
	writer.Write(returnJson)
}
