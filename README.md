# 润措域名资产管理平台

## 部署说明

### Caddy

```Caddyfile
runcuo.com parking.runcuo.com {
        redir https://www.runcuo.com{uri}
}
www.runcuo.com {
        tls 1@5.nu
        root /home/www/runcuo/frontend
        proxy /upload 127.0.0.1:8035 {
                transparent
                keepalive 10
        }
        proxy /api 127.0.0.1:8034 {
                transparent
                keepalive 10
        }
}
localhost:8035 {
        root /home/www/runcuo/upload/
}
* {
        tls {
                ask http://127.0.0.1:8034/allowed
        }
        proxy / 127.0.0.1:8034 {
                transparent
                keepalive 10
        }
        proxy /upload 127.0.0.1:8035 {
                transparent
                keepalive 10
        }
}
```