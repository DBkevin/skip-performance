# 皮肤科绩效管理系统 - 部署指南

## 系统信息

- **前端域名**: clm.xmmylike.com
- **后端端口**: 8111 (本地)
- **数据库**: 远程MySQL (120.25.70.117:3306)
- **管理员账户**: admin / admin123

## Ubuntu服务器部署步骤

### 1. 环境准备

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git nginx certbot python3-certbot-nginx
```

### 2. 安装Go环境

```bash
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
go version
```

### 3. 部署后端

```bash
cd /var/www/backend
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
go mod tidy
go build -o server .
```

### 4. 部署前端

```bash
cd /var/www/frontend
npm install
npm run build
```

### 5. 配置Nginx

```bash
sudo cp /var/www/nginx.conf /etc/nginx/nginx.conf
sudo nginx -t
sudo systemctl restart nginx
sudo systemctl enable nginx
```

### 6. 申请SSL证书

```bash
sudo certbot --nginx -d clm.xmmylike.com
```

### 7. 访问系统

- https://clm.xmmylike.com
- 管理员账户: admin / admin123