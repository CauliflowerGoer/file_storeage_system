package handler

import (
	"encoding/json"
	"file_storage_stystem/meta"
	"file_storage_stystem/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)



// 实现的是返回网址效果



// 用户返回的 response 对象,  用户请求的 request 指针
func UploadHanler(w http.ResponseWriter, r *http.Request)  {
	// 返回上传 html 页面
	if r.Method == "GET" {

		// 读取文件并返回内容
	  data , err := 	ioutil.ReadFile("./static/view/index.html")
	  // 出现错误,
	  if(err != nil){
		  io.WriteString(w, "internel server error")
		  return
	  }
	  // 文件内容写入 response
	  io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录

		file ,  head, err := r.FormFile("file")
		if err != nil {
			fmt.Errorf("%v",err)
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileSha1: "",
			FileName: head.Filename,

			Location: "/Users/hfcb/Desktop/file_storage_stystem/temp/"+head.Filename,
			UploadAt: time.Now().Format("2006-01-02 03:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		defer newFile.Close()
		if err != nil {
			fmt.Errorf("%v",err)
		}



		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// 更新信息

		// 上传  返回字节长度和错误信息
		size, err := io.Copy(newFile,file)

		if(err != nil){
			fmt.Errorf("%v",err)
		}

		fileMeta.FileSize = size/1014

		meta.UpdateFileMeta(fileMeta)
		// 重定向
		http.Redirect(w,r, "/file/upload/succeed",http.StatusFound)
	}
}

// 上传完成
func UploadSucceedHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 获取文件信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	// 解析客户端请求的参数
	r.ParseForm()
	// 获取参数
	filehash :=	r.Form["filehash"][0]
	fmate := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fmate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	// 根据 sha1 获取对应的 fileMeta
	fm := meta.GetFileMeta(fsha1)
	// 读取
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data,err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")

	w.Write(data)
}

// 重命名
func RenameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sh1 := r.Form.Get("filehash")
	newName := r.Form.Get("name")



	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	fm := meta.GetFileMeta(sh1)
	file_path_and_name := fm.Location

	file_paht := []rune(file_path_and_name)[:strings.LastIndex(file_path_and_name,"/")+1]
	os.Rename(file_path_and_name , string(file_paht)+ newName)
	fm.FileName = newName
	meta.UpdateFileMeta(fm)



	data,err := json.Marshal(fm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 文件删除
func FileDelHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form.Get("filehash")


	// 获取 路径, 删除文件
	filePath := meta.GetFileMeta(filehash).Location
	fmt.Println(filePath)
	os.Remove(filePath)


	meta.RemoveFileMeta(filehash)




	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "del finished!")

}