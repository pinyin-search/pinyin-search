package web

import (
	"log"
	"net/http"
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

// Insert insert
func Insert(writer http.ResponseWriter, request *http.Request) {
	data := request.FormValue("data")
	words := jieba.Cut(data, true)
	pys := make([]string, 0)
	for _, word := range words {
		if utf8.RuneCountInString(word) > 1 {
			pys = append(pys, append(pinyin.LazyPinyin(word, pyArgs), strings.Join(pinyin.LazyPinyin(word, pyArgsFirst), ""))...)
		}
	}

	words2 := strings.Fields(data)
	for _, word := range words2 {
		pys = append(pys, strings.Join(regWord.FindAllString(word, -1), ""))
	}

	log.Println(pys)

}
