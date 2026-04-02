# 学习成长管理平台（MVP）设计文档

## 1. 项目概述

**项目名称**：学习成长管理平台（本地运行工具）  
**设计日期**：2026-04-02  
**定位**：面向学生和自学者的轻量级个人效率工具，连接“规划 -> 执行 -> 复盘”完整闭环。  

平台核心不是单纯记录任务，而是把以下场景打通：

1. 今天要学什么
2. 今天实际学了多久
3. 完成了哪些任务
4. 有没有坚持打卡
5. 本周/本月学习状态如何

## 2. 目标用户

1. 大学生
2. 考研/考公/考证用户
3. 自学编程用户
4. 希望管理每日学习计划的人

## 3. MVP 目标与范围

### 3.1 核心目标

首版只解决三件事：

1. 规划：今天学什么
2. 执行：实际学了多久
3. 复盘：最近学得怎么样

### 3.2 首版包含模块

1. 今日任务
2. 学习记录（手动 + 计时器）
3. 每日打卡
4. 统计面板
5. 本地账号体系（单账号）

### 3.3 首版不做内容（明确排除）

1. 云端同步/多设备同步
2. 多账号切换
3. 自动备份与导入导出
4. 上线部署（仅本地运行）

## 4. 产品形态与页面结构

### 4.1 运行形态

1. 第一阶段：本地 Web（浏览器访问 `http://localhost:<port>`）
2. 第二阶段：可封装为桌面应用（保留架构兼容性）

### 4.2 页面结构

1. 登录/注册页
2. 首页统计面板（Dashboard）
3. 今日任务页
4. 学习记录页
5. 每日打卡页
6. 设置页（账号信息、修改密码）

## 5. 功能设计

### 5.1 今日任务模块

**目标**：解决“今天学什么”。

**能力**：

1. 查看今日任务列表
2. 新增任务（标题、优先级、截止日期）
3. 任务编辑/删除
4. 标记完成/取消完成
5. 未完成任务自动顺延到次日

**优先级定义**：

1. `HIGH`
2. `MEDIUM`
3. `LOW`

### 5.2 学习记录模块

**目标**：解决“今天实际学了什么”。

**能力**：

1. 手动录入学习记录（科目、时长、备注）
2. 简单计时器（开始/暂停/继续/结束）
3. 记录编辑/删除
4. 科目可选择已有科目或新增

**记录示例**：

1. 英语，50 分钟，背单词
2. 数学，90 分钟，刷题

### 5.3 每日打卡模块

**目标**：解决“有没有坚持”。

**打卡规则（已确认）**：

1. 当天必须存在至少 1 条有效学习记录
2. 用户必须手动点击“打卡”
3. 同时满足以上两条才算当日完成打卡

**能力**：

1. 一键打卡
2. 防重复打卡（幂等）
3. 展示连续打卡天数

### 5.4 统计面板模块

**目标**：解决“最近表现怎么样”。

**指标卡片**：

1. 今日学习总时长
2. 本周学习总时长（周一开始）
3. 今日完成任务数
4. 当前连续打卡天数

**图表**：

1. 近 7 天学习时长趋势（折线图）
2. 按科目学习时长占比（饼图）

## 6. 技术架构

### 6.1 技术选型

1. 前端：React + TypeScript + Vite
2. 后端：Go 1.24（Gin）
3. ORM：Gorm
4. 数据库：MySQL 8.0.27（本地）
5. 图表：ECharts

### 6.2 架构分层

**前端层**：页面、状态管理、图表展示、表单交互  
**API 层**：RESTful 接口  
**业务层**：任务、学习、打卡、统计规则  
**数据层**：Gorm + MySQL 持久化  

### 6.3 后端模块划分（Go）

1. `internal/auth`：注册、登录、改密、会话校验
2. `internal/tasks`：任务 CRUD、完成状态、顺延逻辑
3. `internal/study`：学习记录与计时器
4. `internal/checkin`：打卡校验、打卡写入、连续天数
5. `internal/stats`：聚合统计查询
6. `internal/subjects`：科目管理
7. `internal/shared`：中间件、错误模型、时间工具、事务工具

## 7. 数据库设计（MySQL 8.0.27）

统一约定：核心业务表都包含 `ext JSON NULL`、`created_at`、`updated_at`，用于后续扩展与兼容升级。

### 7.1 users

1. `id` BIGINT PK AUTO_INCREMENT
2. `username` VARCHAR(64) NOT NULL UNIQUE
3. `password_hash` VARCHAR(255) NOT NULL
4. `ext` JSON NULL
5. `created_at` DATETIME NOT NULL
6. `updated_at` DATETIME NOT NULL

### 7.2 subjects

1. `id` BIGINT PK AUTO_INCREMENT
2. `user_id` BIGINT NOT NULL
3. `name` VARCHAR(64) NOT NULL
4. `color` VARCHAR(16) NULL
5. `ext` JSON NULL
6. `created_at` DATETIME NOT NULL
7. `updated_at` DATETIME NOT NULL

索引：

1. UNIQUE(`user_id`, `name`)

### 7.3 tasks

1. `id` BIGINT PK AUTO_INCREMENT
2. `user_id` BIGINT NOT NULL
3. `title` VARCHAR(255) NOT NULL
4. `priority` ENUM('HIGH','MEDIUM','LOW') NOT NULL DEFAULT 'MEDIUM'
5. `due_date` DATE NULL
6. `plan_date` DATE NOT NULL
7. `status` ENUM('PENDING','DONE') NOT NULL DEFAULT 'PENDING'
8. `completed_at` DATETIME NULL
9. `carry_count` INT NOT NULL DEFAULT 0
10. `ext` JSON NULL
11. `created_at` DATETIME NOT NULL
12. `updated_at` DATETIME NOT NULL

索引：

1. INDEX(`user_id`, `plan_date`, `status`)
2. INDEX(`user_id`, `due_date`)

### 7.4 study_sessions

1. `id` BIGINT PK AUTO_INCREMENT
2. `user_id` BIGINT NOT NULL
3. `subject_id` BIGINT NOT NULL
4. `record_type` ENUM('MANUAL','TIMER') NOT NULL
5. `start_at` DATETIME NOT NULL
6. `end_at` DATETIME NOT NULL
7. `duration_minutes` INT NOT NULL
8. `note` VARCHAR(1000) NULL
9. `ext` JSON NULL
10. `created_at` DATETIME NOT NULL
11. `updated_at` DATETIME NOT NULL

索引：

1. INDEX(`user_id`, `start_at`)
2. INDEX(`user_id`, `subject_id`, `start_at`)

### 7.5 daily_checkins

1. `id` BIGINT PK AUTO_INCREMENT
2. `user_id` BIGINT NOT NULL
3. `checkin_date` DATE NOT NULL
4. `checked_at` DATETIME NOT NULL
5. `ext` JSON NULL
6. `created_at` DATETIME NOT NULL
7. `updated_at` DATETIME NOT NULL

索引：

1. UNIQUE(`user_id`, `checkin_date`)

### 7.6 timer_states

1. `id` BIGINT PK AUTO_INCREMENT
2. `user_id` BIGINT NOT NULL
3. `status` ENUM('IDLE','RUNNING','PAUSED') NOT NULL DEFAULT 'IDLE'
4. `subject_id` BIGINT NULL
5. `started_at` DATETIME NULL
6. `last_resumed_at` DATETIME NULL
7. `paused_seconds` INT NOT NULL DEFAULT 0
8. `draft_note` VARCHAR(1000) NULL
9. `ext` JSON NULL
10. `created_at` DATETIME NOT NULL
11. `updated_at` DATETIME NOT NULL

索引：

1. UNIQUE(`user_id`)

## 8. 核心业务规则

### 8.1 账号规则

1. 首次启动必须注册本地账号（用户名 + 密码）
2. 首版仅支持单账号使用
3. 数据模型保留 `user_id` 以支持后续多账号扩展

### 8.2 任务顺延规则

1. 每日首次登录后执行顺延流程
2. 将昨天 `PENDING` 的任务直接更新为今天的 `plan_date`
3. 每次顺延 `carry_count + 1`
4. 顺延流程必须使用事务
5. 同一天重复触发顺延时不应重复处理（幂等）

### 8.3 学习记录规则

1. `duration_minutes` 必须大于 0
2. `end_at` 必须晚于 `start_at`
3. 手动录入和计时器记录统一写入 `study_sessions`
4. 删除记录后统计结果必须实时变化

### 8.4 打卡规则

1. 当天存在学习记录且手动点击打卡才可成功
2. 当天重复打卡返回“已打卡”，不重复写入
3. 连续打卡仅基于 `daily_checkins.checkin_date` 连续自然日计算

### 8.5 统计规则

1. 今日时长：当天学习时长总和
2. 本周时长：以周一为周起点的自然周总和
3. 今日完成任务数：当天标记为 `DONE` 的任务数
4. 科目占比：按科目汇总学习时长占比
5. 7 天趋势：最近 7 天每日学习总时长（含 0 时长日期）

## 9. API 设计（MVP）

### 9.1 认证

1. `POST /api/auth/register`
2. `POST /api/auth/login`
3. `POST /api/auth/change-password`
4. `GET /api/auth/me`

### 9.2 科目

1. `GET /api/subjects`
2. `POST /api/subjects`
3. `PUT /api/subjects/:id`
4. `DELETE /api/subjects/:id`

### 9.3 任务

1. `GET /api/tasks/today`
2. `POST /api/tasks`
3. `PUT /api/tasks/:id`
4. `PATCH /api/tasks/:id/status`
5. `DELETE /api/tasks/:id`

### 9.4 学习记录/计时器

1. `GET /api/study-sessions?date=YYYY-MM-DD`
2. `POST /api/study-sessions/manual`
3. `PUT /api/study-sessions/:id`
4. `DELETE /api/study-sessions/:id`
5. `POST /api/timer/start`
6. `POST /api/timer/pause`
7. `POST /api/timer/resume`
8. `POST /api/timer/stop`
9. `GET /api/timer/state`

### 9.5 打卡与统计

1. `POST /api/checkins/today`
2. `GET /api/checkins/streak`
3. `GET /api/stats/overview`
4. `GET /api/stats/weekly-trend`
5. `GET /api/stats/subject-distribution`

## 10. 异常处理与稳定性

1. MySQL 不可用：返回统一错误码并在前端提示“数据库连接失败，请检查本地 MySQL 服务”
2. 参数非法：统一返回 400，包含字段级错误信息
3. 业务冲突（重复打卡/重复科目名）：返回 409
4. 顺延任务、计时器结束写记录等关键操作使用事务
5. 计时器异常中断后可从 `timer_states` 恢复状态

## 11. 测试策略

### 11.1 单元测试（Go）

1. 连续打卡天数计算
2. 周区间计算（周一为起点）
3. 任务顺延幂等性
4. 统计聚合函数准确性

### 11.2 接口测试

1. 认证流程（注册/登录/鉴权）
2. 任务 CRUD 与状态切换
3. 学习记录 CRUD 与计时器流程
4. 打卡前置条件校验与幂等
5. 统计接口返回格式与数值正确性

### 11.3 端到端最小链路

1. 注册账号
2. 登录
3. 新增今日任务
4. 新增学习记录（手动或计时器）
5. 手动打卡
6. 验证统计面板数据变化

## 12. 里程碑建议（实现顺序）

1. 项目脚手架与基础设施（前后端、数据库连接、迁移）
2. 认证与会话
3. 今日任务模块
4. 学习记录与计时器
5. 每日打卡
6. 统计面板与图表
7. 联调与验收测试


