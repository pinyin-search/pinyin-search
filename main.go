package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"pinyin-suggestion/web"
	"time"
)

// 监听地址
var listen = flag.String("l", ":9876", "监听地址")

func main() {
	log.Println("监听", *listen, "...")

	http.HandleFunc("/insert", web.Insert)
	err := http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Println("启动端口发生异常, 请检查端口是否被占用", err)
		time.Sleep(time.Minute)
		os.Exit(1)
	}

}
