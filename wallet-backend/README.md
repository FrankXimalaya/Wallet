# Wallet Backend

以太坊钱包后端服务

## 功能特性

- 用户管理（注册/登录）
- 钱包地址生成和管理
- ETH 转账功能
- 交易记录查询
- 区块链监听（充值检测）

## 项目结构

```
wallet-backend/
├── cmd/server/          # 应用启动入口
├── internal/            # 内部实现
│   ├── config/         # 配置管理
│   ├── eth/            # 以太坊相关功能
│   ├── http/           # HTTP API
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   └── service/        # 业务逻辑层
└── scripts/            # 脚本和配置
```

## 快速开始

1. 安装依赖
```bash
go mod download
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和以太坊节点
```

3. 启动服务
```bash
# 使用 docker-compose 启动数据库
cd scripts && docker-compose up -d

# 运行服务
go run cmd/server/main.go
```

## API 接口

### 用户相关
- POST /api/users/register - 用户注册
- POST /api/users/login - 用户登录

### 钱包相关
- POST /api/wallet/create - 创建钱包地址
- GET /api/wallet/:address - 查询钱包余额

### 交易相关
- POST /api/tx/transfer - 发起转账
- GET /api/tx/:hash - 查询交易详情
- GET /api/tx/list/:address - 查询地址交易记录

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM
- PostgreSQL
- go-ethereum
