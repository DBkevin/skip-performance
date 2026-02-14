# 皮肤科绩效管理系统

一个基于 Go + Vue3 的皮肤科医院绩效管理平台。

## 技术栈

### 后端
- **Go** - 主要开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL 8.0** - 数据库
- **JWT** - 认证授权

### 前端
- **Vue3** - 前端框架
- **Element Plus** - UI组件库
- **Vite** - 构建工具
- **Axios** - HTTP客户端

### 数据库
- MySQL 8.0
- 支持软删除
- 完整的外键关联

## 功能特性

### 核心模块
1. **顾客管理** - 初诊/复诊/再消费分类
2. **员工管理** - 医生/护士/咨询师/管理员角色
3. **项目管理** - 美容/医疗项目维护
4. **就诊管理** - 单据录入、业绩分配
5. **回访记录** - 护士回访统计
6. **业绩报表** - 多维度统计分析

### 业绩分配逻辑
- **主操医生**: 总金额 × (1 - 协同比例总和)
- **协同医生**: 按设定比例分配
- **护士**: 固定比例(如5%)

### 权限控制
- 管理员: 可管理员工/项目
- 普通用户: 只能录入和查看

## 快速开始

### 后端部署

```bash
cd backend

# 设置Go代理
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 下载依赖
go mod tidy

# 编译
go build -o server .

# 运行
./server
```

后端默认运行在 `:8080`

### 前端部署

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

前端默认运行在 `:3000`

### 数据库配置

修改 `backend/config/database.go`:
```go
dsn := "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

## 项目结构

```
skip-performance/
├── backend/          # Go后端
│   ├── config/       # 数据库配置
│   ├── controllers/  # API控制器
│   ├── middleware/   # 中间件
│   ├── models/       # 数据库模型
│   └── routes/       # 路由配置
├── frontend/         # Vue3前端
│   ├── src/
│   │   ├── api/      # API封装
│   │   ├── components/
│   │   ├── router/   # 路由配置
│   │   ├── utils/    # 工具函数
│   │   └── views/    # 页面组件
│   └── package.json
└── deploy.sh         # 部署脚本
```

## API接口

### 认证
- `POST /api/login` - 登录

### 顾客管理
- `GET /api/customers` - 顾客列表
- `POST /api/customers` - 创建顾客
- `PUT /api/customers/:id` - 更新顾客
- `DELETE /api/customers/:id` - 删除顾客

### 员工管理 (需管理员权限)
- `GET /api/employees` - 员工列表
- `POST /api/employees` - 创建员工
- `PUT /api/employees/:id` - 更新员工
- `DELETE /api/employees/:id` - 删除员工

### 项目管理 (需管理员权限)
- `GET /api/projects` - 项目列表
- `POST /api/projects` - 创建项目
- `PUT /api/projects/:id` - 更新项目
- `DELETE /api/projects/:id` - 删除项目

### 就诊管理
- `GET /api/visits` - 就诊列表
- `POST /api/visits` - 创建就诊
- `PUT /api/visits/:id` - 更新就诊
- `DELETE /api/visits/:id` - 删除就诊

### 业绩报表
- `GET /api/reports/performance` - 业绩统计

## 开发计划

- [x] 基础架构搭建
- [x] 数据库模型设计
- [x] API接口开发
- [x] 前端页面开发
- [ ] 数据导入导出
- [ ] 图表可视化
- [ ] 移动端适配
- [ ] Docker部署

## License

MIT
