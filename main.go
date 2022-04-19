package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"pinyin-search/search"
	"pinyin-search/web"
	"time"
)

// 监听地址
var listen = flag.String("l", ":7701", "监听地址")
var authorization = flag.String("auth", "", "认证, 留空不进行认证")
var meiliHost = flag.String("meiliHost", "127.0.0.1:7700", "meilisearch Host")
var meiliKey = flag.String("meiliKey", "", "meilisearch APIKey")

func main() {
	flag.Parse()
	if _, err := net.ResolveTCPAddr("tcp", *listen); err != nil {
		log.Fatalf("解析监听地址异常，%s", err)
	}

	os.Setenv(web.AUTHORIZATION_ENV, *authorization)
	os.Setenv(search.MEILISEARCH_HOST_ENV, *meiliHost)
	os.Setenv(search.MEILISEARCH_APIKEY_ENV, *meiliKey)

	// init search
	if *meiliHost != "" {
		search.MySearch = &search.MeiliSearch{}
		search.MySearch.Init()
	}

	http.HandleFunc("/add", web.Auth(web.Add))
	http.HandleFunc("/suggestion", web.Auth(web.Suggestion))

	log.Println("监听", *listen, "...")

	err := http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Println("启动端口发生异常, 请检查端口是否被占用", err)
		time.Sleep(time.Minute)
		os.Exit(1)
	}

}
