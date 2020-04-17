package main

import (
	"fmt"
	"net/http"

	"file_storage_stystem/handler"
)

func main() {
	// 路径与处理对象
	//mux := http.NewServeMux()
	http.HandleFunc("/file/upload", handler.UploadHanler)
	http.HandleFunc("/file/upload/succeed",handler.UploadSucceedHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/rename", handler.RenameHandler)

	http.HandleFunc("/file/remove", handler.FileDelHandler)
	//mux.Handle("/file" ,http.FileServer(http.Dir(".")) )


	// 监听端口
	err := http.ListenAndServe(":8000",nil)
	if err != nil {
		fmt.Printf("failded to start server ,%v",err)
	}
}
