# xwc1125-pkg

## 简介

`xwc1125-pkg` golang通用库

## 功能

- [x] freecache
- [x] captcha
- [x] copier: copy struct to struct
- [x] mysql: gorm & xorm
- [x] es
- [x] kafka
- [x] redis
- [x] sqlite
- [x] email
- [x] di 基于反射实现依赖注入
- [x] goinject: fx依赖注入
- [x] hashring 哈希环
- [x] jwtauth
- [x] ldap
- [x] otp 二次认证
- [x] pprof二次封装
- [x] rbac
- [x] secure
    - [x] checkurl 校验URL的合法性，用于防止跳转漏洞、SSRF漏洞
    - [x] ddm 动态数据掩码,防止敏感数据暴露
    - [x] ipfilter ip过滤
    - [x] password 对密码进行混淆加密
    - [x] safecurl 安全的http请求，防止SSRF
    - [x] shift 位移算法
    - [x] ssrf 判断url是否会触发ssrf
- [x] snowflake 雪花算法
- [ ] ssh
- [x] util
- [x] validator
- [x] version 版本管理
- [x] watemark 添加水印

## 使用

获取包 `go get github.com/xwc1125/xwc1125-pkg`

## 证书

`xwc1125-pkg` 的源码允许用户在遵循 [Apache 2.0 开源证书](LICENSE) 规则的前提下使用。

## 版权

Copyright@2022 xwc1125

![xwc1125](./logo.png)
