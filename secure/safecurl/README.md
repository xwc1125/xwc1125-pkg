## scurl

本函数类用于安全地做HTTP请求，用于防止 SSRF(服务器端请求伪造)漏洞。

### 简介

| 函数                   | 说明                            |
|----------------------|-------------------------------|
| secapi.NewSafeClient | 创建安全的HTTP请求客户端（推荐）            |
| secapi.SafeCurl      | 安全地发起HTTP GET请求，并返回response对象 |
| secapi.GetSafeURL    | 检查传入的URL是否指向内网资源              |

