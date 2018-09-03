# 润措域名资产管理平台

## 部署说明

### Caddy

```Caddyfile
localhost:8080{
    /static {
        root /static #后端 theme 目录
    }
    proxy /api 127.0.0.1:8081/api
    root / #前端 dist 目录
}
```