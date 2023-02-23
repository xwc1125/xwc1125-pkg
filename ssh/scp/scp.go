// Package scp
//
// @author: xwc1125
package scp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"github.com/xwc1125/xwc1125-pkg/ssh/sftpx"
	"golang.org/x/crypto/ssh"
)

func toUnixPath(path string) string {
	return filepath.Clean(path)
}

// RZ 上传文件
// @param h 远程服务器信息
// @param localPath 本地地址
// @param remotePath 远程地址
func RZ(ctx context.Context, client *ssh.Client, localPath string, remotePath string) error {
	return Scp(ctx, true, client, localPath, remotePath)
}

// SZ 下载文件
// @param h 远程服务器信息
// @param remotePath 远程地址
// @param localPath 本地地址
func SZ(ctx context.Context, client *ssh.Client, remotePath, localPath string) error {
	return Scp(ctx, false, client, localPath, remotePath)
}

// Scp 进行scp处理
// @param localToRemote 是否本地复制到远程服务器
// @param h 远程服务器信息
// @param localPath 本地地址
// @param remotePath 远程地址
func Scp(ctx context.Context, localToRemote bool, client *ssh.Client, localPath string, remotePath string) error {
	c, err := sftpx.NewSFTP(ctx, client)
	if err != nil {
		return err
	}
	defer c.Close()
	return scpCopy(localToRemote, c.SftpClient(), remotePath, localPath)
}

// scpCopy 进行scp处理
func scpCopy(localToRemote bool, c *sftp.Client, remote, local string) error {
	var (
		info os.FileInfo
		err  error
	)
	if localToRemote {
		info, err = os.Lstat(local)
	} else {
		info, err = c.Lstat(toUnixPath(remote))
	}
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return scpCopyLink(localToRemote, c, local, remote)
	}
	if info.IsDir() {
		return scpCopyD(localToRemote, c, remote, local)
	}
	return scpCopyF(localToRemote, c, remote, local, info)
}

// scpCopyLink 复制软链
func scpCopyLink(localToRemote bool, c *sftp.Client, remote, local string) error {
	var (
		realLocal  = local
		realRemote = remote
		err        error
	)

	if localToRemote {
		realLocal, err = os.Readlink(local)
	} else {
		realRemote, err = c.ReadLink(toUnixPath(remote))
	}
	if err != nil {
		return err
	}
	return scpCopy(localToRemote, c, realRemote, realLocal)
}

// scpCopyD 复制文件夹
func scpCopyD(localToRemote bool, c *sftp.Client, remote, local string) error {
	if localToRemote {
		contents, err := ioutil.ReadDir(local)
		if err != nil {
			return fmt.Errorf("ioutil read local dir failed %s", err)
		}
		for _, content := range contents {
			cdL, csR := filepath.Join(local, content.Name()), filepath.Join(remote, content.Name())
			if err := scpCopy(localToRemote, c, csR, cdL); err != nil {
				return fmt.Errorf("%w %s %s", err, cdL, csR)
			}
		}
		return nil
	}
	contents, err := c.ReadDir(toUnixPath(remote))
	if err != nil {
		return fmt.Errorf("ioutil read scp remote dir failed %s", err)
	}
	for _, info := range contents {
		cdL, csR := filepath.Join(local, info.Name()), filepath.Join(remote, info.Name())
		// 本地创建文件夹
		err := os.MkdirAll(filepath.Dir(cdL), info.Mode())
		if err != nil {
			return fmt.Errorf("os local sub mkdir all failed,%s", err)
		}
		csR = toUnixPath(csR)
		if err := scpCopy(localToRemote, c, csR, cdL); err != nil {
			return fmt.Errorf("dir walk remote:%s, local:%s, %s", csR, cdL, err)
		}
	}
	return nil
}

// scpCopyF 复制文件
func scpCopyF(localToRemote bool, c *sftp.Client, remote, local string, info os.FileInfo) error {
	if localToRemote {
		localFile, err := os.Open(local)
		if err != nil {
			return fmt.Errorf("BrowserOpen local file failed %w", err)
		}
		defer localFile.Close()
		err = c.MkdirAll(toUnixPath(filepath.Dir(remote)))
		if err != nil {
			return fmt.Errorf("scp mkdir all failed %w", err)
		}
		remoteFile, err := c.Create(toUnixPath(remote))
		if err != nil {
			return fmt.Errorf("create remote file failed %s:%s", remote, err)
		}
		defer remoteFile.Close()
		err = c.Chmod(remoteFile.Name(), info.Mode())
		if err != nil {
			return fmt.Errorf("scp chmod failed %s", err)
		}
		_, err = io.Copy(remoteFile, localFile)
		if err != nil {
			return fmt.Errorf("io copy failed %s", err)
		}
		return nil
	}
	rFile, err := c.Open(toUnixPath(remote))
	if err != nil {
		return fmt.Errorf("BrowserOpen scp remote file failed %s", err)
	}
	defer rFile.Close()

	lFile, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("os create local file failed:%s %s", local, err)
	}
	defer lFile.Close()

	size, err := io.Copy(lFile, rFile)
	if err != nil {
		return fmt.Errorf("io copy remote to local failed.size:%d %s", size, err)
	}

	err = os.Chmod(lFile.Name(), info.Mode())
	if err != nil {
		return fmt.Errorf("os local chmod failed %s", err)
	}
	return nil
}
