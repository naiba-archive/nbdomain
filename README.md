# 润措域名资产管理平台

## 部署说明

### Caddy

```Caddyfile
:80 {
   redir https://{hostname}{uri} 
}
:443 {
    proxy / 127.0.0.1:8081
}
runcuo.com{
    /static {
        root /static #米表主题目录
    }
    /upload {
        root /upload #上传文件目录
    }
    proxy /api 127.0.0.1:8081/api
    root / #前端 dist 目录
}
```