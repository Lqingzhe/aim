## 项目说明
### 总览
本项目是一个生产级的分布式即时通讯系统后端，采用微服务架构，支持单聊、群聊、消息推送、文件传输、AI 助手等功能。附带面向开发者的简单前端页面。
#### 后端技术栈
| 类别     | 技术            | 版本      | 用途             |
|--------|---------------|---------|----------------|
| 语言     | go            | 1.26.2  | 主要开发语言         |
| Web框架  | gin           | v1.12.0 | HTTP路由和中间件     |
| RPC框架  | Kitex         | v0.16.1 | 微服务通信          |
| 服务发现   | Nacos         | v1.1.6  | 服务注册与发现        |
| 关系型数据库 | Mysql(GORM)   | v1.31.1 | 持久化储存          |
| 缓存     | Redis         | v9.18.0 | 缓存、限流以及部分持久化储存 |
| 文档数据库  | MongoDB       | v2.5.1  | 消息以及ai上下文储存    |
| 消息队列   | Kafka(sarama) | v1.48.1 | 异步消息           |
| 日志     | zap           | v1.27.1 | 结构化日志          |
| A框架    | Eino          | -       | AI Agent编排     |

### 核心功能
***1、用户板块***
- 用户注册/登录（JWT Token 认证）
- AccessToken + RefreshToken 双 Token 机制
- 用户信息管理（头像、昵称、个性签名、生日）
- 用户备注功能
- 多设备登录管理（单设备/全部设备退出）
***2、群聊模块***
- 群聊创建/解散
- 群信息修改/搜索
- 入群申请（申请→审核→同意/拒绝）
- 成员管理（踢出、转让群主、设置/撤销管理员）
- 禁言功能（时长+原因）
***3、私聊模块***
- 好友申请（发送→同意/拒绝）
- 私聊会话管理
- 删除好友  
***4、消息模块***
- 文本消息发送/接收
- 图片/文件/语音发送
- 消息撤回
- 历史消息加载（分页）
- 新消息提醒
- 离线消息存储
- 消息已读状态  
***5、实时通信***
- WebSocket 长连接
- 心跳保活
- 消息实时推送
- 断线自动重连
- Token 过期自动刷新  
***6、AI助手***
- 用户可配置自己的 AI 参数（模型、API地址、角色、提示词）
- 群聊中 @bot 触发 AI 对话
- AI Agent 支持 Tool 调用
- 由AI自身管理用户相关的用户画像
- 对话历史管理
- 上下文长度控制（包括对话轮数以及上下文总长度）
##### 详情点击下方按钮
[gateway](./app/api/README.md)  
[user-service](./app/userservice/README.md)  
[group-service](./app/groupservice/README.md)  
[message-service](./app/messageservice/README.md)  
[file-service](./app/fileservice/README.md)  
[ai-service](./app/aiservice/README.md)  
[my-error](./pkg/error/README.md)  
[my--log](./pkg/log/README.md)
## 项目结构
```
IM项目/
├── app/                           # 所有微服务
│   ├── gateway/                   # API网关服务
│   │   ├── api/                   # HTTP路由配置
│   │   ├── config/                # 配置管理
│   │   ├── handler/               # 请求处理层
│   │   ├── middleware/            # 中间件层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   ├── main.go
│   │   ├── templates/             # HTML模板
│   │   │   ├── login.html         # 登录页面
│   │   │   └── chat.html          # 聊天页面
│   │   │
│   │   └── static/                # 静态资源
│   │       ├── css/
│   │       │   └── chat.css       # 聊天样式
│   │       └── js/                # 前端JS文件
│   │
│   ├── userservice/               # 用户服务
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   └── main.go
│   │
│   ├── groupservice/              # 群组服务
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   └── main.go
│   │
│   ├── messageservice/            # 消息服务
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   └── main.go
│   │
│   ├── fileservice/               # 文件服务
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   └── main.go
│   │
│   └─ aiservice/                  # ai服务
│       ├── config/                # 配置管理
│       ├── dao/                   # 数据访问层
│       ├── handler/               # RPC处理层
│       ├── model/                 # 数据模型
│       ├── service/               # 业务逻辑层
│       ├── agent/                 # angent管理
│       └── main.go
│
├── commonmodel/                   # 公共数据模型
│   ├── common_config_model.go     # 配置模型
│   ├── const.go                   # 常量定义
│   ├── db_config_model.go         # 数据库配置模型
│   ├── interface.go               # 通用接口
│   ├── kafka_message_model.go     # Kafka消息模型
│   └── redis_lua.go               # Redis Lua脚本
│
├── kitex_gen/                     # Kitex生成的RPC代码
│   ├── kitexcommonmodel/          # 公共RPC模型
│   ├── kitexfileservice/          # 文件服务RPC
│   ├── kitexgroupservice/         # 群组服务RPC
│   ├── kitexmessageservice/       # 消息服务RPC
│   └── kitexuserservice/          # 用户服务RPC
│
├── pkg/                           # 公共包
│   ├── commonconfig/              # 配置加载
│   │   ├── config.go              # 配置解析
│   │   ├── init_db.go             # 数据库初始化
│   │   ├── init_kafka.go          # Kafka初始化
│   │   └── init_service.go        # 服务注册发现
│   │
│   ├── error/                     # 错误处理
│   │   ├── error.go               # 错误定义
│   │   ├── error_code.go          # 错误码
│   │   └── db_error.go            # 数据库错误处理
│   │
│   ├── id/                        # ID生成器
│   │   └── snow.go                # Snowflake雪花算法
│   │
│   ├── log/                       # 日志系统
│   │   ├── log.go                 # 日志初始化
│   │   └── set_log.go             # Gin日志中间件
│   │
│   └── commondao/                 # 通用DAO层
│       └── dao.go                 # L1缓存实现
│
├── tool/                          # 工具函数
│   ├── db_tool.go                 # 数据库工具
│   ├── kafka_tool.go              # Kafka工具
│   └── tool.go                    # 通用工具
│
│
├── config.yaml                    # 配置文件
└── go.mod                         # Go模块依赖
```
## 接口说明
### 基础信息
- **响应格式**：JSON
- **认证方式**:
    - 需要登录的接口需要在 Header 中添加 `Authorization: Bearer {access_token}`
    - 同时需要添加 `X-Device-ID` Header 标识设备
    - websocket的路由仅需{access_token}

### 公开接口（无需登录）

| 方法   | 路径                    | 用途             | 请求参数                                                              | 响应示例                                                                                                      |
|------|-----------------------|----------------|-------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------|
| GET  | `/`                   | 根路径重定向         | 无                                                                 | 重定向到 `/login`                                                                                             |
| GET  | `/login`              | 登录/注册页面        | 无                                                                 | HTML 页面                                                                                                   |
| GET  | `/ping`               | 健康检查           | 无                                                                 | `{"code":0,"message":"pong"}`                                                                             |
| POST | `/user/register`      | 用户注册           | `{"password":"string"}`                                           | `{"code":0,"message":"success","data":{"user_info":{"user_id":"string"}}}`                                |
| POST | `/user/login`         | 用户登录           | `{"user_id":"string","password":"string"}` + Header `X-Device-ID` | `{"code":0,"message":"success","data":{"token_info":{"access_token":"string","refresh_token":"string"}}}` |
| POST | `/user/refresh-token` | 刷新 AccessToken | `{"refresh_token":"string"}` + Header `X-Device-ID`               | `{"code":0,"message":"success","data":{"token_info":{"access_token":"string","refresh_token":"string"}}}` |

---

### 需要登陆的接口
以下接口需要在Header中添加：  
`Authorization: Bearer {access_token} `  
`X-Device-ID: {device_id} `  

#### 用户模块 (`/user`)

| 方法   | 路径                          | 用途        | 请求参数                                                                                                   | 响应示例                                                                                                                                                                          |
|------|-----------------------------|-----------|--------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| POST | `/user/logout-all-device`   | 退出所有设备    | 无                                                                                                      | `{"code":0,"message":"Logout Success"}`                                                                                                                                       |
| POST | `/user/logout-a-device`     | 退出当前设备    | 无                                                                                                      | `{"code":0,"message":"Logout Success"}`                                                                                                                                       |
| POST | `/user/get-user-info`       | 获取当前用户信息  | 无                                                                                                      | `{"code":0,"message":"success","data":{"user_info":{"UserInfo":{...},"RemarkInfos":[...]}}}`                                                                                  |
| POST | `/user/get-other-user-info` | 获取其他用户信息  | `{"goal_user_id":"string"}`                                                                            | `{"code":0,"message":"success","data":{"user_info":{"user_name":"string","introduction":"string","birthday_year":0,"birthday_month":0,"birthday_day":0,"is_connect":false}}}` |
| POST | `/user/update-user-info`    | 更新用户信息    | `{"user_name":"string","introduction":"string","birthday_year":0,"birthday_month":0,"birthday_day":0}` | `{"code":0,"message":"success"}`                                                                                                                                              |
| POST | `/user/remark`              | 设置/删除用户备注 | `{"goal_user_id":"string","nick_name":"string"}` (nick_name 为空表示删除备注)                                  | `{"code":0,"message":"success"}`                                                                                                                                              |

---

#### 群组模块 (`/group`)

##### 群聊管理

| 方法   | 路径                                   | 用途            | 请求参数                                                 | 响应示例                                                                                                                                           |
|------|--------------------------------------|---------------|------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| POST | `/group/create-group`                | 创建群聊          | `{"group_name":"string"}`                            | `{"code":0,"message":"success","data":{"group_info":{"group_id":"string"}}}`                                                                   |
| POST | `/group/delete-group`                | 解散群聊（仅群主）     | `{"group_id":"string"}`                              | `{"code":0,"message":"success"}`                                                                                                               |
| POST | `/group/leave-group`                 | 退出群聊          | `{"group_id":"string"}`                              | `{"code":0,"message":"success"}`                                                                                                               |
| POST | `/group/get-group-info`              | 获取群聊信息        | `{"group_id":"string"}`                              | `{"code":0,"message":"success","data":{"group_info":{"group_id":"string","group_name":"string"}}}`                                             |
| POST | `/group/change-group-info`           | 修改群聊信息        | `{"group_id":"string","group_name":"string"}`        | `{"code":0,"message":"success"}`                                                                                                               |
| POST | `/group/search-group`                | 搜索群聊          | `{"group_name":"string"}`                            | `{"code":0,"message":"success","data":{"group_info":{"group_id_list":[]}}}`                                                                    |
| POST | `/group/get-group-info-with-user`    | 获取用户在群中的信息    | `{"group_id":"string"}`                              | `{"code":0,"message":"success","data":{"group_info":{"group_id":"string","group_remark_name":"string","group_role":"string"}}}`                |
| POST | `/group/update-group-info-with-user` | 修改用户在群中的备注    | `{"group_id":"string","group_remark_name":"string"}` | `{"code":0,"message":"success"}`                                                                                                               |
| POST | `/group/get-group-and-session-id`    | 获取用户所有群聊和私聊会话 | 无                                                    | `{"code":0,"message":"success","data":{"session_info":{"session_id_list":[],"user_of_session_id_list":[]},"group_info":{"group_id_list":[]}}}` |

##### 群申请管理

| 方法   | 路径                            | 用途              | 请求参数                                            | 响应示例                                                                        |
|------|-------------------------------|-----------------|-------------------------------------------------|-----------------------------------------------------------------------------|
| POST | `/group/set-group-apply`      | 申请加入群聊          | `{"group_id":"string"}`                         | `{"code":0,"message":"success"}`                                            |
| POST | `/group/get-group-apply-list` | 获取群申请列表（群主/管理员） | `{"group_id":"string"}`                         | `{"code":0,"message":"success","data":{"group_info":{"group_id_list":[]}}}` |
| POST | `/group/agree-group-apply`    | 同意入群申请          | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}`                                            |
| POST | `/group/refuse-group-apply`   | 拒绝入群申请          | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}`                                            |

##### 群成员管理

| 方法   | 路径                             | 用途    | 请求参数                                            | 响应示例                             |
|------|--------------------------------|-------|-------------------------------------------------|----------------------------------|
| POST | `/group/transform-group-owner` | 转让群主  | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/set-manager`           | 设置管理员 | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/revoke-manager`        | 撤销管理员 | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/kick-out-group`        | 踢出群成员 | `{"group_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success"}` |

##### 禁言管理

| 方法   | 路径                    | 用途   | 请求参数                                                                                         | 响应示例                             |
|------|-----------------------|------|----------------------------------------------------------------------------------------------|----------------------------------|
| POST | `/group/set-mute`     | 禁言成员 | `{"group_id":"string","goal_user_id":"string","mute_time_seconds":0,"mute_reason":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/release-mute` | 解除禁言 | `{"group_id":"string","goal_user_id":"string"}`                                              | `{"code":0,"message":"success"}` |

##### 已读状态

| 方法   | 路径                           | 用途             | 请求参数                    | 响应示例                                                                                    |
|------|------------------------------|----------------|-------------------------|-----------------------------------------------------------------------------------------|
| POST | `/group/get-last-visit-time` | 获取群成员最后访问时间    | `{"group_id":"string"}` | `{"code":0,"message":"success","data":{"group_info":{"last_visit_time":{"用户ID":时间戳}}}}` |
| POST | `/group/set-last-visit-time` | 更新用户在群中的最后访问时间 | `{"group_id":"string"}` | `{"code":0,"message":"success"}`                                                        |

---

#### 好友/私聊模块 (`/group`)

| 方法 | 路径 | 用途 | 请求参数 | 响应示例 |
|------|------|------|----------|----------|
| POST | `/group/apply-for-friend` | 发送好友申请 | `{"goal_user_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/get-friend-apply-list` | 获取好友申请列表 | 无 | `{"code":0,"message":"success","data":{"session_info":{"apply_user_list":[]}}}` |
| POST | `/group/refuse-friend-apply` | 拒绝好友申请 | `{"goal_user_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/creat-session` | 同意好友申请并创建会话 | `{"goal_user_id":"string"}` | `{"code":0,"message":"success","data":{"session_info":{"session_id":"string"}}}` |
| POST | `/group/delete-session` | 删除私聊会话（删除好友） | `{"session_id":"string"}` | `{"code":0,"message":"success"}` |
| POST | `/group/get-friend-last-visit-time` | 获取好友最后访问时间 | `{"session_id":"string","goal_user_id":"string"}` | `{"code":0,"message":"success","data":{"session_info":{"goal_user_id":"string","last_visit_time":0}}}` |

---
#### 消息模块 (`/message`)

| 方法   | 路径                           | 用途         | 请求参数                                                                 | 响应示例                                                                             |
|------|------------------------------|------------|----------------------------------------------------------------------|----------------------------------------------------------------------------------|
| POST | `/message/send-message`      | 发送文本消息     | `{"group_id":"string","message_content":"string"}`                   | `{"code":0,"message":"success","data":{"message_info":{"message_id":"string"}}}` |
| POST | `/message/send-file`         | 发送文件       | `multipart/form-data`:<br>- `group_id`<br>- `file_name`<br>- `file`  | `{"code":0,"message":"success","data":{"message_info":{"message_id":"string"}}}` |
| POST | `/message/send-voice`        | 发送语音       | `multipart/form-data`:<br>- `group_id`<br>- `voice_time`<br>- `file` | `{"code":0,"message":"success","data":{"message_info":{"message_id":"string"}}}` |
| POST | `/message/send-picture`      | 发送图片       | `multipart/form-data`:<br>- `group_id`<br>- `picture`                | `{"code":0,"message":"success","data":{"message_info":{"message_id":"string"}}}` |
| POST | `/message/withdraw-message`  | 撤回消息       | `{"group_id":"string","message_id":"string"}`                        | `{"code":0,"message":"success"}`                                                 |
| POST | `/message/get-message-list`  | 获取历史消息     | `{"group_id":"string","start_time_second":0,"end_time_second":0}`    | `{"code":0,"message":"success","data":{"message_info":{"message_list":[...]}}}`  |
| POST | `/message/get-new-message`   | 获取新消息      | `{"group_id":"string"}`                                              | `{"code":0,"message":"success","data":{"message_list":[...]}}`                   |
| POST | `/message/get-file-content`  | 下载文件/图片/语音 | `{"group_id":"string","message_id":"string"}`                        | 文件流                                                                              |
| POST | `/message/send-group-notice` | 发送群公告      | `{"group_id":"string","message_content":"string"}`                   | `{"code":0,"message":"success"}`                                                 |

---

#### AI 服务模块 (`/ai`)

| 方法   | 路径                        | 用途            | 请求参数                                                                                               | 响应示例                                                                                                                                                   |
|------|---------------------------|---------------|----------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| POST | `/ai/delete-chat-context` | 删除用户的AI对话历史记录 | 无                                                                                                  | `{"code":0,"message":"success"}`                                                                                                                       |
| POST | `/ai/get-ai-config`       | 获取用户的AI配置     | 无                                                                                                  | `{"code":0,"message":"success","data":{"ai_config":{"model_name":"string","base_url":"string","api_key":"string","role":"string","prompt":"string"}}}` |
| POST | `/ai/update-ai-config`    | 更新/创建用户的AI配置  | `{"model_name":"string","base_url":"string","api_key":"string","role":"string","prompt":"string"}` | `{"code":0,"message":"success"}`                                                                                                                       |
| POST | `/ai/delete-ai-config`    | 删除用户的AI配置     | 无                                                                                                  | `{"code":0,"message":"success"}`                                                                                                                       |

##### AI 配置字段说明

| 字段           | 类型     | 说明                             |
|--------------|--------|--------------------------------|
| `model_name` | string | AI模型名称（如 gpt-4、qwen-max、glm-4） |
| `base_url`   | string | API服务地址                        |
| `api_key`    | string | API密钥                          |
| `role`       | string | AI角色设定（如：助手、编程专家）              |
| `prompt`     | string | 系统提示词                          |

##### 使用流程

1. **配置AI**：用户需要先调用 `/ai/update-ai-config` 配置自己的AI参数
2. **@bot对话**：在群聊中发送 `@bot 你的问题` 即可与AI对话
3. **删除历史**：调用 `/ai/delete-chat-context` 清除对话历史记录

---


### WebSocket 连接

| 路径    | 用途     | 连接参数                    | 说明                         |
|-------|--------|-------------------------|----------------------------|
| `/ws` | 实时消息推送 | `?token={access_token}` | 建立 WebSocket 连接后，可接收实时消息推送 |

#### WebSocket 消息格式

| 消息类型            | 说明        | 示例                                              |
|-----------------|-----------|-------------------------------------------------|
| `ping` / `pong` | 心跳检测      | 客户端发送 `{"type":"ping"}`，服务端回复 `{"type":"pong"}` |
| `logout`        | 退出登录通知    | `{"type":"logout"}`                             |
| `group_message` | 群聊新消息通知   | 包含消息内容的 JSON                                    |
| 系统通知            | 好友申请、群申请等 | 包含 message_code 的 JSON                          |