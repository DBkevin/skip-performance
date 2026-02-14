#!/bin/bash

# 皮肤科绩效管理系统部署脚本

set -e

echo "=== 开始部署 ==="

# 检查Go是否已安装
if ! command -v go &> /dev/null; then
    echo "错误: Go 未安装"
    exit 1
fi

# 检查Node.js是否已安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js 未安装"
    exit 1
fi

# 设置部署目录
DEPLOY_DIR="/var/www/skip"
WORKSPACE_DIR="/home/lee/.openclaw/workspace/skip"

# 后端部署
echo "=== 部署后端 ==="
cd "$WORKSPACE_DIR/backend"
go mod tidy
go build -o bin/server .

# 创建部署目录（如果有权限）
if [ -d "$DEPLOY_DIR" ]; then
    echo "部署目标目录已存在"
else
    echo "警告: $DEPLOY_DIR 不存在，尝试创建..."
    mkdir -p "$DEPLOY_DIR" 2>/dev/null || echo "无权限创建 $DEPLOY_DIR，将部署到用户目录"
    DEPLOY_DIR="$HOME/www/skip"
    mkdir -p "$DEPLOY_DIR"
fi

# 复制后端文件
mkdir -p "$DEPLOY_DIR/backend"
cp -r bin "$DEPLOY_DIR/backend/" 2>/dev/null || true
cp -r config "$DEPLOY_DIR/backend/" 2>/dev/null || true

echo "后端部署完成"

# 前端部署
echo "=== 部署前端 ==="
cd "$WORKSPACE_DIR/frontend"
npm install
npm run build

# 复制前端文件
mkdir -p "$DEPLOY_DIR/frontend"
cp -r dist/* "$DEPLOY_DIR/frontend/"

echo "前端部署完成"

# 创建启动脚本
cat > "$DEPLOY_DIR/start.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")/backend"
nohup ./bin/server > server.log 2>&1 &
echo $! > server.pid
echo "服务已启动，PID: $(cat server.pid)"
EOF

chmod +x "$DEPLOY_DIR/start.sh"

cat > "$DEPLOY_DIR/stop.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")/backend"
if [ -f server.pid ]; then
    kill $(cat server.pid) 2>/dev/null || true
    rm server.pid
    echo "服务已停止"
else
    echo "服务未运行"
fi
EOF

chmod +x "$DEPLOY_DIR/stop.sh"

echo "=== 部署完成 ==="
echo "部署目录: $DEPLOY_DIR"
echo "启动命令: $DEPLOY_DIR/start.sh"
echo "停止命令: $DEPLOY_DIR/stop.sh"
echo "日志文件: $DEPLOY_DIR/backend/server.log"
