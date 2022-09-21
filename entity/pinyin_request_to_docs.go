package entity

import (
	"encoding/json"
	"errors"
	"fmt"
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
	seg.LoadDictEmbed()
}

func (req *PinYinRequest) GetDocs() (docs []Doc, err error) {

	// 校验
	if req.IndexName == "" || req.DataId == "" || req.Data == "" {
		r, _ := json.Marshal(req)
		return docs, errors.New("异常的数据: " + string(r))
	}

	words := seg.Cut(req.Data, true)
	indexes := make(map[string]string)

	// 全拼字母
	quanpin := pinyin.LazyPinyin(req.Data, pyArgs)
	if len(quanpin) > 0 {
		indexes[strings.Join(quanpin, "")] = req.Data
	}

	// 全拼第一个字母
	firstLetterAll := strings.Join(pinyin.LazyPinyin(req.Data, pyArgsFirst), "")
	if firstLetterAll != "" {
		indexes[firstLetterAll] = req.Data
	}

	// 分词
	for _, word := range words {
		if utf8.RuneCountInString(word) > 1 {
			// 首字母
			firstLetter := strings.Join(pinyin.LazyPinyin(word, pyArgsFirst), "")
			if firstLetter != "" {
				indexes[firstLetter] = word
			} else {
				// 英文单词 或 数字
				indexes[word] = word
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
	for _, word := range strings.Fields(req.Data) {
		if word != "" {
			wordsEnglish = wordsEnglish + " " + strings.Join(regWord.FindAllString(word, -1), "")
		}
	}
	wordsEnglish = strings.TrimSpace(wordsEnglish)
	if wordsEnglish != "" {
		indexes[wordsEnglish] = wordsEnglish
	}

	// 原始数据
	indexes[req.Data] = req.Data

	idx := 0
	for k, v := range indexes {
		// id 为 dataId_index
		docs = append(docs, Doc{Id: fmt.Sprintf("%s_%d", req.DataId, idx), Key: k, Value: v})
		idx = idx + 1
	}

	return docs, nil
}
