package meta

import "fmt"

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}
// 新增/更新 文件元 信息
func UpdateFileMeta(fmeta FileMeta) {
	fmt.Println(fmeta)
	fileMetas[fmeta.FileSha1] = fmeta
}

// 通过 sha1 值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 从 map 中删除数据
func RemoveFileMeta(filesha1 string)  {
	delete(fileMetas,filesha1)

}