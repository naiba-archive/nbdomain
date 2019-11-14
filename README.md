# 日落米表托管

## 部署说明

### Caddy

```Caddyfile
riluo.cn parking.riluo.cn {
	redir https://www.riluo.cn{uri}
}
www.riluo.cn {
	tls 1@5.nu
	root /home/www/runcuo/frontend
	proxy /api 127.0.0.1:8034 {
				timeout 5m
        		keepalive 50
                transparent
                keepalive 10
        }
}
www.riluo.cn/upload/ {
	root /home/www/runcuo/upload/
}
www.riluo.cn/static/ {
        root /home/www/runcuo/theme/static/
}

:80 {
        proxy / 127.0.0.1:8034 {
                timeout 5m
                keepalive 50
                transparent
                keepalive 10
        }
}
:443 {
        tls {
                ask http://127.0.0.1:8034/allowed
        }
        proxy / 127.0.0.1:8034 {
        		timeout 5m
        		keepalive 50
                transparent
                keepalive 10
        }
}
```