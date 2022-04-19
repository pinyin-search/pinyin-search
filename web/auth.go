package web

import (
	"encoding/json"
	"net/http"
	"os"
	"pinyin-search/entity"
)

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

const AUTHORIZATION_ENV = "AUTHORIZATION_PINYIN"

// Auth auth
func Auth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		authorization := os.Getenv(AUTHORIZATION_ENV)
		if authorization != "" && auth != authorization {
			// 执行被装饰的函数
			returnjson, _ := json.Marshal(entity.Result{Success: false, Msg: "认证失败"})
			w.Write(returnjson)
			return
		}
		// 执行被装饰的函数
		f(w, r)
	}
}
