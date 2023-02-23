// Package sftpx
//
// @author: xwc1125
package sftpx

import "mime/multipart"

type Mode int

const (
	Mode_ReadDir  Mode = iota + 1 // 读取列表
	Mode_Mkdir                    // 创建文件夹
	Mode_Rename                   // 修改名称
	Mode_Remove                   // 删除
	Mode_Download                 // 下载
	Mode_Upload                   // 上传
)

type OptParams struct {
	Mode     Mode           `json:"mode"`                // 操作类型
	Path     string         `json:"path,omitempty"`      // 路径
	IsDir    bool           `json:"is_dir,omitempty"`    // 路径是否为文件夹
	OldName  string         `json:"old_name"`            // 旧文件名
	NewName  string         `json:"new_name"`            // 新文件名
	FileName string         `json:"file_name,omitempty"` // 文件上传时的文件名
	File     multipart.File `json:"file,omitempty"`      // 文件上传时的文件内容
}

// FileInfo 文件信息
type FileInfo struct {
	Name    string `json:"name"`    // 文件名称
	Path    string `json:"path"`    // 路径
	Size    int64  `json:"size"`    // 文件大小
	Mode    string `json:"mode"`    // 权限
	ModTime string `json:"modTime"` // 修改时间
	IsDir   bool   `json:"isDir"`   // 是否是文件夹
}

// FileInfos 文件对象
type FileInfos struct {
	List []FileInfo `json:"list"`
	Dir  string     `json:"dir"`
}
