# 学习成长管理平台（本地 MVP）使用文档

本文档是完整的中文启动手册，覆盖：
- 本地 MySQL 启动（Docker / 本机 MySQL 两种方式）
- 后端启动
- 前端启动
- `TestMVPFlow` 跑成真实 PASS（非 SKIP）
- 常见问题排查

## 1. 项目结构与运行端口

- 后端：`backend/`（Go 1.24 + Gin + Gorm）
- 前端：`frontend/`（Vite + React + TypeScript）
- 数据库：MySQL 8.0.27

默认端口：
- 后端 API：`http://localhost:8080`
- 前端开发服务：`http://localhost:5173`
- MySQL：`127.0.0.1:3306`

## 2. 环境要求

- Go `1.24+`
- Node.js `18+`
- npm `10+`（推荐）
- MySQL 8.0.27（可用 Docker 或本机服务）

## 3. 数据库连接信息（固定）

项目默认数据库参数：
- 用户名：`root`
- 密码：`010511`
- 数据库：`learning_growth`
- DSN：

```txt
root:010511@tcp(127.0.0.1:3306)/learning_growth?charset=utf8mb4&parseTime=True&loc=Local
```

## 4. 启动 MySQL（推荐 Docker）

在项目根目录执行：

```powershell
docker compose up -d mysql
```

检查容器状态：

```powershell
docker compose ps
```

检查健康状态（可选）：

```powershell
docker inspect --format='{{json .State.Health.Status}}' learning-growth-mysql
```

### 4.1 如果你使用的是本机 MySQL（非 Docker）

确保服务已启动，并创建数据库：

```powershell
mysql -uroot -p010511 -e "CREATE DATABASE IF NOT EXISTS learning_growth CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

验证数据库存在：

```powershell
mysql -uroot -p010511 -e "SHOW DATABASES LIKE 'learning_growth';"
```

## 5. 启动后端服务

在 `backend/` 下准备环境变量（Windows PowerShell）：

```powershell
Set-Location backend
Copy-Item .env.example .env -Force
```

编辑 `.env`，确认至少包含：

```env
APP_PORT=8080
MYSQL_DSN=root:010511@tcp(127.0.0.1:3306)/learning_growth?parseTime=true&loc=Local&charset=utf8mb4
JWT_SECRET=local-dev-secret
```

启动后端：

```powershell
go run ./cmd/server
```

健康检查：

```powershell
curl http://localhost:8080/api/health
```

预期返回：`{"status":"ok"}`

## 6. 启动前端服务

在新终端中执行：

```powershell
Set-Location frontend
npm install
npm run dev
```

打开：`http://localhost:5173`

如果后端不是默认地址，可设置：

```powershell
$env:VITE_API_BASE_URL="http://localhost:8080"
npm run dev
```

## 7. 跑测试（含 TestMVPFlow）

### 7.1 后端集成测试（重点）

进入 `backend/`：

```powershell
Set-Location backend
$env:MYSQL_DSN='root:010511@tcp(127.0.0.1:3306)/learning_growth?charset=utf8mb4&parseTime=True&loc=Local&timeout=2s&readTimeout=2s&writeTimeout=2s'
go test ./internal/integration -run TestMVPFlow -v -count=1 -timeout=180s
```

如果 MySQL 可用且 `learning_growth` 存在，预期是：
- `--- PASS: TestMVPFlow`
- `ok   learning-growth-platform/internal/integration`

### 7.2 后端全量测试

```powershell
go test ./... -count=1 -timeout=180s
```

### 7.3 前端测试与构建

进入 `frontend/`：

```powershell
Set-Location ../frontend
npm test -- --run
npm run build
```

## 8. TestMVPFlow 覆盖的业务流程

`TestMVPFlow` 会模拟完整 MVP 主链路：
1. 注册
2. 登录
3. 创建科目
4. 创建今日任务（DONE）
5. 创建今日学习记录
6. 今日打卡
7. 获取 `/api/stats/overview` 并校验关键指标

## 9. 常见问题排查

### 9.1 `Unknown database 'learning_growth'`

原因：MySQL 已启动，但数据库未创建。  
处理：

```powershell
mysql -uroot -p010511 -e "CREATE DATABASE IF NOT EXISTS learning_growth CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 9.2 `TestMVPFlow` 显示 `SKIP`

原因：测试探测到 MySQL 依赖不可用（连接失败/数据库不存在）。  
处理顺序：
1. 确认 3306 端口可用
2. 确认 `root/010511` 可登录
3. 确认 `learning_growth` 已创建
4. 重新执行 `go test ./internal/integration -run TestMVPFlow -v`

### 9.3 前端 `vite/vitest` 报 `spawn EPERM`

常见于受限沙箱/权限环境。  
建议：
1. 在普通本机终端执行 `npm test -- --run` 和 `npm run build`
2. 确认杀软未阻止 `node`/`esbuild`
3. 删除 `node_modules` 后重装：`npm install`

## 10. 一次性完整启动命令清单（Windows PowerShell）

```powershell
# 1) 项目根目录
Set-Location D:\ai_coding\learning-management-system\.worktrees\codex-learning-growth-mvp

# 2) 启动 MySQL（Docker 方式）
docker compose up -d mysql

# 3) 后端
Set-Location backend
Copy-Item .env.example .env -Force
go run ./cmd/server

# 4) 新终端启动前端
Set-Location D:\ai_coding\learning-management-system\.worktrees\codex-learning-growth-mvp\frontend
npm install
npm run dev
```
