// Package sftpx
//
// @author: xwc1125
package sftpx

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/chain5j/logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"github.com/xwc1125/xwc1125-pkg/types/response"
	"golang.org/x/crypto/ssh"
)

var (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxPacket = 1 << 15
)

// SFTP sftp对象
type SFTP struct {
	log    logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	sftpClient *sftp.Client
}

// NewSFTP 创建sftp
func NewSFTP(rootCtx context.Context, client *ssh.Client) (*SFTP, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	// 此时获取了sshClient，下面使用sshClient构建sftpClient
	sftpClient, err := sftp.NewClient(client, sftp.MaxPacket(maxPacket))
	if err != nil {
		return nil, err
	}

	return &SFTP{
		log:    logger.Log("sftp"),
		ctx:    ctx,
		cancel: cancel,

		sftpClient: sftpClient,
	}, nil
}

func (s *SFTP) ServeWs(wsConn *websocket.Conn) error {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if wsConn != nil {
					if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						wsConn.Close()
						return
					}
				}
			case <-s.ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			break
		}
		if wsConn != nil {
			messageType, message, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					logger.Error("ws close err", "err", err)
				}
				break
			}
			if messageType == websocket.TextMessage {
				code := 0
				msg := ""
				fileInfos, err := s.ReadDir(string(message))
				if err != nil {
					code = 1
					msg = err.Error()
				}

				b, _ := json.Marshal(response.Response{
					Code: code,
					Msg:  msg,
					Data: fileInfos,
				})
				if err := wsConn.WriteMessage(websocket.TextMessage, b); err != nil {
					logger.Error("ws write message err", "err", err)
					break
				}
			}
		}
	}
	return nil
}

func getFileInfos(dirPath string, fileInfos []os.FileInfo) *FileInfos {
	var fileList = make([]FileInfo, 0)
	for _, f := range fileInfos {
		if f.Mode()&os.ModeSymlink != 0 {
			continue
		}
		fileList = append(fileList, FileInfo{
			Name:    f.Name(),
			Path:    path.Join(dirPath, f.Name()),
			Size:    f.Size(),
			Mode:    f.Mode().String(),
			ModTime: f.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   f.IsDir(),
		})
	}
	return &FileInfos{
		List: fileList,
		Dir:  dirPath,
	}
}

// Close 关闭
func (s *SFTP) Close() error {
	s.cancel()
	if s.sftpClient != nil {
		return s.sftpClient.Close()
	}
	return nil
}

// ReadDir 读取文件夹
func (s *SFTP) ReadDir(dirPath string) (*FileInfos, error) {
	if dirPath == "" {
		if dirPath == "" {
			dir, err := s.sftpClient.Getwd()
			if err != nil {
				return nil, err
			}
			dirPath = dir
		}
	}
	fileInfos, err := s.sftpClient.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return getFileInfos(dirPath, fileInfos), nil
}

// Mkdir 创建文件夹
func (s *SFTP) Mkdir(dirPath string) error {
	err := s.sftpClient.MkdirAll(dirPath)
	if err != nil {
		return err
	}
	return nil
}

// Rename 重命名
func (s *SFTP) Rename(o, n string) error {
	err := s.sftpClient.Rename(o, n)
	if err != nil {
		return err
	}
	return nil
}

// DownloadDir 下载文件/夹
func (s *SFTP) DownloadDir(fullPath string) ([]byte, error) {
	sftpClient := s.sftpClient

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	err := zipAddFiles(w, sftpClient, fullPath, "/")
	if err != nil {
		return nil, err
	}
	// 确保写入无错误
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// zipAddFiles 压缩所有的文件
func zipAddFiles(w *zip.Writer, sftpC *sftp.Client, basePath, baseInZip string) error {
	// 打开文件夹
	files, err := sftpC.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("sftp read dir %s failed:%s", basePath, err)
	}

	// 循环所有的文件
	for _, file := range files {
		thisFilePath := basePath + "/" + file.Name()
		if file.IsDir() {
			// 如果是文件夹，递归处理
			err := zipAddFiles(w, sftpC, thisFilePath, baseInZip+file.Name()+"/")
			if err != nil {
				return fmt.Errorf("zip add files %s failed:%s", thisFilePath, err)
			}
		} else {
			// 文件，直接打开
			dat, err := sftpC.Open(thisFilePath)
			if err != nil {
				return fmt.Errorf("sftp open file %s err:%s", thisFilePath, err)
			}
			// 归档
			zipElePath := baseInZip + file.Name()
			f, err := w.Create(zipElePath)
			if err != nil {
				return fmt.Errorf("zip create path %s err:%s", zipElePath, err)
			}
			b, err := ioutil.ReadAll(dat)
			if err != nil {
				return fmt.Errorf("ioutil read all failed %s", err)
			}
			_, err = f.Write(b)
			if err != nil {
				return fmt.Errorf("zip write data err:%s", err)
			}
		}
	}
	return nil
}

// DownloadFile 下载文件
func (s *SFTP) DownloadFile(fullPath string) ([]byte, error) {
	f, err := s.sftpClient.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// Upload 上传
func (s *SFTP) Upload(desDir string, fileName string, srcFile multipart.File) error {
	if desDir == "$HOME" {
		wd, err := s.sftpClient.Getwd()
		if err != nil {
			return err
		}
		desDir = wd
	}

	// 创建文件
	dstFile, err := s.sftpClient.Create(path.Join(desDir, fileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 写入内容
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return err
	}
	return nil
}

// Remove 删除
func (s *SFTP) Remove(fullPath string, isDir bool) error {
	if fullPath == "/" || fullPath == "$HOME" {
		return errors.New("can't delete / or $HOME dir")
	}
	if isDir {
		return s.removeNonemptyDirectory(fullPath)
	}
	return s.sftpClient.Remove(fullPath)
}

// removeNonemptyDirectory 删除非空目录.
// sftp协议不允许删除非空目录，因此需要遍历文件树以有序地删除文件和目录
func (s *SFTP) removeNonemptyDirectory(path string) error {
	list, err := s.sftpClient.ReadDir(path)
	if err != nil {
		return err
	}
	// 遍历文件树
	for i, cur := range list {
		newPath := filepath.Join(path, list[i].Name())
		if cur.IsDir() {
			if err := s.removeNonemptyDirectory(newPath); err != nil {
				return err
			}
		} else {
			if err := s.sftpClient.Remove(newPath); err != nil {
				return err
			}
		}
	}
	// 删除当前文件夹，因为它已经为空
	return s.sftpClient.RemoveDirectory(path)
}

func (s *SFTP) SftpClient() *sftp.Client {
	return s.sftpClient
}

func (s *SFTP) HandleOpt(params OptParams) (interface{}, error) {
	switch params.Mode {
	case Mode_ReadDir:
		fileInfos, err := s.ReadDir(params.Path)
		if err != nil {
			return nil, err
		}
		return fileInfos, nil
	case Mode_Mkdir:
		err := s.Mkdir(params.Path)
		return nil, err
	case Mode_Rename:
		err := s.Rename(params.OldName, params.NewName)
		return nil, err
	case Mode_Remove:
		err := s.Remove(params.Path, params.IsDir)
		return nil, err
	case Mode_Download:
		if params.IsDir {
			bs, err := s.DownloadDir(params.Path)
			if err != nil {
				return nil, err
			}
			return bs, err
		} else {
			bs, err := s.DownloadFile(params.Path)
			if err != nil {
				return nil, err
			}
			return bs, err
		}
	case Mode_Upload:
		err := s.Upload(params.Path, params.FileName, params.File)
		return nil, err
	}
	return nil, errors.New("unsupported the mode")
}
