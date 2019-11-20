# 日落米表托管

## 部署说明

### Caddy

```caddyfile
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

### 宝塔部署

1. 复制二进制、生成数据目录 `mkdir -p data/upload/logo`
2. 克隆主题目录 `git clone https://github.com/naiba/nbdomain-theme.git`
3. 使用 `systemd` 守护后端进程
4. 创建 MySQL 数据库
5. 启动后端进程
6. 创建管理面板站点，复制静态文件到站点目录
7. 配置站点的反向代理 

   ```conf
   location / {
        # 用于配合 browserHistory使用
        try_files $uri $uri/ /index.html;
    }
    location /api/ {
        expires 12h;
        if ($request_uri ~* "(php|jsp|cgi|asp|aspx)")
        {
             expires 0;
        }
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host localhost:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header REMOTE-HOST $remote_addr;

        #持久化连接相关配置
        #proxy_connect_timeout 30s;
        #proxy_read_timeout 86400s;
        #proxy_send_timeout 30s;
        #proxy_http_version 1.1;
        #proxy_set_header Upgrade $http_upgrade;
        #proxy_set_header Connection "upgrade";
        add_header X-Cache $upstream_cache_status;

        #Set Nginx Cache
        add_header Cache-Control no-cache;
    }
    location /static/ {
        expires 12h;
        if ($request_uri ~* "(php|jsp|cgi|asp|aspx)")
        {
             expires 0;
        }
        proxy_pass http://localhost:8080/static/;
        proxy_set_header Host localhost:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header REMOTE-HOST $remote_addr;

        #持久化连接相关配置
        #proxy_connect_timeout 30s;
        #proxy_read_timeout 86400s;
        #proxy_send_timeout 30s;
        #proxy_http_version 1.1;
        #proxy_set_header Upgrade $http_upgrade;
        #proxy_set_header Connection "upgrade";
        add_header X-Cache $upstream_cache_status;

        #Set Nginx Cache
        add_header Cache-Control no-cache;
    }
    location /upload/ {
        expires 12h;
        if ($request_uri ~* "(php|jsp|cgi|asp|aspx)")
        {
             expires 0;
        }
        proxy_pass http://localhost:8080/upload/;
        proxy_set_header Host localhost:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header REMOTE-HOST $remote_addr;

        #持久化连接相关配置
        #proxy_connect_timeout 30s;
        #proxy_read_timeout 86400s;
        #proxy_send_timeout 30s;
        #proxy_http_version 1.1;
        #proxy_set_header Upgrade $http_upgrade;
        #proxy_set_header Connection "upgrade";
        add_header X-Cache $upstream_cache_status;

        #Set Nginx Cache
        add_header Cache-Control no-cache;
    }
   ```

8. 修改默认站点的反向代理

   ```conf
   location / {
        expires 12h;
        if ($request_uri ~* "(php|jsp|cgi|asp|aspx)")
        {
             expires 0;
        }
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header REMOTE-HOST $remote_addr;

        #持久化连接相关配置
        #proxy_connect_timeout 30s;
        #proxy_read_timeout 86400s;
        #proxy_send_timeout 30s;
        #proxy_http_version 1.1;
        #proxy_set_header Upgrade $http_upgrade;
        #proxy_set_header Connection "upgrade";
        add_header X-Cache $upstream_cache_status;

        #Set Nginx Cache
        add_header Cache-Control no-cache;
    }
   ```
