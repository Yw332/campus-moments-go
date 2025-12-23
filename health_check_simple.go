package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 简单的健康检查服务
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"code":200,"message":"success","data":{"status":"ok","service":"Simple Health Check"}}`)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"code":200,"message":"Campus Moments API is running","data":{}}`)
	})

	fmt.Println("🚀 简单健康检查服务启动...")
	fmt.Println("📡 地址: http://0.0.0.0:8080")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
