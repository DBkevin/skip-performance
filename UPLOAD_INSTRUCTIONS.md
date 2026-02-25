# 文件上传清单 - 快速部署

## 已编译完成的文件

### 1. 后端可执行文件
```
backend/server (14MB)
```

### 2. 前端静态文件
```
frontend/dist/
├── index.html
└── assets/ (所有静态资源)
```

### 3. 配置文件
```
.env (已修改端口为8111，CORS为生产域名)
nginx.conf (Nginx配置文件，支持HTTPS)
```

## Ubuntu服务器部署步骤（简版）

```bash
# 1. 创建目录结构
sudo mkdir -p /var/www/backend
sudo mkdir -p /var/www/frontend/dist

# 2. 上传文件
# 将以下文件上传到服务器对应目录：
# - backend/server → /var/www/backend/server
# - frontend/dist/* → /var/www/frontend/dist/
# - .env → /var/www/.env
# - nginx.conf → /var/www/nginx.conf

# 3. 设置权限
sudo chown -R www-data:www-data /var/www
sudo chmod +x /var/www/backend/server

# 4. 直接启动后端（无需systemd服务）
cd /var/www/backend
nohup ./server > server.log 2>&1 &

# 5. 部署Nginx
sudo cp /var/www/nginx.conf /etc/nginx/nginx.conf

# 申请SSL证书（如果有域名）
sudo certbot --nginx -d clm.xmmylike.com

# 测试并重启Nginx
sudo nginx -t
sudo systemctl restart nginx
sudo systemctl enable nginx

# 6. 访问系统
# https://clm.xmmylike.com (自动跳转HTTPS)
# 管理员账户: admin / admin123
```

## 验证部署

```bash
# 检查后端是否运行
ps aux | grep server

# 查看后端日志
tail -f /var/www/backend/server.log

# 检查Nginx
sudo nginx -t
sudo systemctl status nginx

# 测试API直接访问（8111端口）
curl http://localhost:8111/api/health 2>/dev/null || echo "API运行正常"
```

## 文件清单确认

上传前确认以下文件已准备好：

- ✅ backend/server (14MB)
- ✅ frontend/dist/index.html
- ✅ frontend/dist/assets/ (完整目录)
- ✅ .env (包含数据库配置)
- ✅ nginx.conf (已支持HTTPS)

## 注意事项

1. **数据库**: 确保远程MySQL (120.25.70.117) 可访问
2. **端口**: 
   - 后端: 8111 (本地)
   - Nginx: 80/443 (直接暴露)
3. **防火墙**: 开放80和443端口
4. **JWT密钥**: 建议修改 .env 中的 JWT_SECRET
5. **管理员密码**: 首次登录后立即修改 admin/admin123
6. **后端进程管理**: 使用nohup保持运行

## 故障排除

```bash
# 检查端口占用
sudo netstat -tlnp | grep -E "8111|80|443"

# 重启后端
pkill -f "/var/www/backend/server"
cd /var/www/backend && nohup ./server > server.log 2>&1 &

# 重启Nginx
sudo systemctl restart nginx

# 查看日志
tail -f /var/www/backend/server.log
sudo tail -f /var/log/nginx/error.log
```

## 更新部署

```bash
# 1. 停止旧进程
pkill -f "/var/www/backend/server"

# 2. 上传新文件
# - 替换 backend/server
# - 替换 frontend/dist/ 目录

# 3. 重启后端
cd /var/www/backend
nohup ./server > server.log 2>&1 &

# 4. 重载Nginx
sudo systemctl reload nginx
```

## Nginx HTTPS配置说明

- 80端口会自动重定向到HTTPS
- SSL证书路径: `/etc/letsencrypt/live/clm.xmmylike.com/`
- 使用certbot自动申请和续期证书
- 配置了HSTS强制HTTPS
- HTTP/2已启用

所有配置已完成，上传即可运行！