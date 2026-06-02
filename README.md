```
IM项目/
├── app/                           # 所有微服务
│   ├── gateway/                   # API网关服务 (端口 8080)
│   │   ├── api/                   # HTTP路由配置
│   │   ├── config/                # 配置管理
│   │   ├── handler/               # 请求处理层
│   │   ├── middleware/            # 中间件层
│   │   ├── model/                 # 数据模型
│   │   ├── service/               # 业务逻辑层
│   │   ├── templates/             # HTML模板
│   │   │   ├── login.html         # 登录页面
│   │   │   └── chat.html          # 聊天页面
│   │   │
│   │   └── static/                # 静态资源
│   │       ├── css/
│   │       │   └── chat.css       # 聊天样式
│   │       └── js/                # 前端JS文件
│   │
│   ├── userservice/               # 用户服务 (端口 8889)
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   └── service/               # 业务逻辑层
│   │
│   ├── groupservice/              # 群组服务 (端口 8890)
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   └── service/               # 业务逻辑层
│   │
│   ├── messageservice/            # 消息服务 (端口 8891)
│   │   ├── config/                # 配置管理
│   │   ├── dao/                   # 数据访问层
│   │   ├── handler/               # RPC处理层
│   │   ├── model/                 # 数据模型
│   │   └── service/               # 业务逻辑层
│   │
│   └── fileservice/               # 文件服务 (端口 8892)
│       ├── config/                # 配置管理
│       ├── dao/                   # 数据访问层
│       ├── handler/               # RPC处理层
│       ├── model/                 # 数据模型
│       └── service/               # 业务逻辑层
│
├── commonmodel/                   # 公共数据模型
│   ├── common_config_model.go    # 配置模型
│   ├── const.go                   # 常量定义
│   ├── db_config_model.go        # 数据库配置模型
│   ├── interface.go               # 通用接口
│   ├── kafka_message_model.go    # Kafka消息模型
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