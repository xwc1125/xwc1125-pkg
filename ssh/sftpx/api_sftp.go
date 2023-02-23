package sftpx

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"github.com/xwc1125/xwc1125-pkg/types/response"
	"golang.org/x/crypto/ssh"
)

type SftpApi struct {
}

// SftpFileInfos 文件夹列表
func (ap SftpApi) SftpFileInfos(c contextx.Context, sshClient *ssh.Client) {
	dirPath := c.Query("path")
	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()
	items, err := client.ReadDir(dirPath)
	if checkErr(c, err) {
		return
	}
	response.OkData(c, items)
}

// SftpRm 删除文件夹
func (ap SftpApi) SftpRm(c contextx.Context, sshClient *ssh.Client) {
	dirPath := c.Query("path")
	isDirStr := c.Query("is_dir")

	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()
	var isDir = false
	if strings.EqualFold(isDirStr, "true") || strings.EqualFold(isDirStr, "1") {
		isDir = true
	}
	err = client.Remove(dirPath, isDir)
	if checkErr(c, err) {
		return
	}
	response.OkDefault(c)
}

// SftpUpload 上传文件
func (ap SftpApi) SftpUpload(c contextx.Context, sshClient *ssh.Client) {
	desDir := c.Query("path")
	formFile, err := c.FormFile("file")
	if checkErr(c, err) {
		return
	}

	srcFile, err := formFile.Open()
	if checkErr(c, err) {
		return
	}
	defer srcFile.Close()

	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()

	err = client.Upload(desDir, formFile.Filename, srcFile)
	if checkErr(c, err) {
		return
	}
	response.OkDefault(c)
}

// SftpDownloadFile 下载文件
func (ap SftpApi) SftpDownloadFile(c contextx.Context, sshClient *ssh.Client) {
	path := c.Query("path")

	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()

	bs, err := client.DownloadFile(path)
	if checkErr(c, err) {
		return
	}
	fn := filepath.Base(path)
	returnFile(c, bs, "octet-stream", fmt.Sprintf(`attachment; filename="%s"`, fn))

}

func returnFile(c contextx.Context, data []byte, ct, cd string) {
	// buff := bytes.NewBuffer(data)
	// c.DataFromReader(200, int64(buff.Len()), ct, buff, map[string]string{"Content-Disposition": cd})
}

// SftpDownloadDir 下载文件夹
func (ap SftpApi) SftpDownloadDir(c contextx.Context, sshClient *ssh.Client) {
	path := c.Query("path")
	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()

	bs, err := client.DownloadDir(path)
	if checkErr(c, err) {
		return
	}
	fn := strings.ReplaceAll(path, "/", "_") + ".zip"
	returnFile(c, bs, "application/zip", fmt.Sprintf(`attachment; filename="%s"`, fn))
}

// SftpMkdir 创建文件
func (ap SftpApi) SftpMkdir(c contextx.Context, sshClient *ssh.Client) {
	path := c.Query("path")
	client, err := NewSFTP(context.Background(), sshClient)
	if checkErr(c, err) {
		return
	}
	defer client.Close()
	err = client.Mkdir(path)
	if checkErr(c, err) {
		return
	}
	response.OkDefault(c)
}

func checkErr(w contextx.Context, err error) bool {
	flag := err != nil
	if flag {
		response.Fail(w, response.StatusBadRequest, err.Error())
	}
	return flag
}
