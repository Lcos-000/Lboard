server/                          # 项目根目录
├── go.mod                       # Go 模块定义
├── go.sum                       # Go 依赖校验文件
├── cmd/                         # # 应用程序入口点
│   ├── whiteboard-server/       #   主服务入口
│   │   └── main.go             #     服务启动入口文件
│   └── tools/                   #   辅助工具集
│       ├── wal-replay/          #     WAL 日志回放工具（用于数据恢复/审计）
│       │   └── main.go
│       └── loadgen/             #     负载生成工具（用于压力测试）
│           └── main.go
├── configs/                     # # 配置文件目录
│   ├── config.yaml              #   默认配置文件
│   └── config.local.yaml        #   本地开发环境配置（覆盖默认）
├── internal/                    # # 内部应用逻辑（不对外暴露）
│   ├── app/                     #   应用组装与启动
│   │   ├── app.go              #     应用主结构体定义
│   │   ├── wire.go             #     Wire 依赖注入配置
│   │   └── lifecycle.go        #     应用生命周期管理（启动/停止）
│   ├── config/                  #   配置读取与管理
│   │   └── config.go           #     配置结构体与加载逻辑
│   ├── http/                    #   HTTP 服务层
│   │   ├── router.go           #     路由注册
│   │   ├── middleware/          #     HTTP 中间件
│   │   │   ├── auth.go         #       认证中间件
│   │   │   ├── cors.go         #       跨域中间件
│   │   │   ├── logging.go      #       请求日志中间件
│   │   │   └── recovery.go     #       panic 恢复中间件
│   │   └── handlers/            #     HTTP 请求处理器
│   │       ├── auth_handler.go  #       认证相关接口处理
│   │       ├── room_handler.go  #       房间相关接口处理
│   │       ├── file_handler.go  #       文件相关接口处理
│   │       └── health_handler.go #      健康检查接口
│   ├── websocket/               #   WebSocket 连接管理
│   │   ├── gateway.go          #     WebSocket 网关（连接路由/鉴权）
│   │   ├── connection.go       #     连接抽象（读写管理）
│   │   ├── hub.go              #     连接池（连接注册/注销/广播）
│   │   ├── message.go          #     消息结构体定义
│   │   ├── writer.go           #     消息写入器（协程安全）
│   │   ├── reader.go           #     消息读取器
│   │   ├── heartbeat.go        #     心跳检测
│   │   └── backpressure.go     #     背压控制（防止内存溢出）
│   ├── room/                    #   房间管理（基于 Actor 模型）
│   │   ├── manager.go          #     房间管理器（创建/查找/销毁房间）
│   │   ├── actor.go            #     房间 Actor 主循环（事件驱动）
│   │   ├── event.go            #     房间事件定义
│   │   ├── state.go            #     房间状态管理（画布数据）
│   │   ├── broadcast.go        #     消息广播（向房间内成员推送）
│   │   ├── lifecycle.go        #     房间生命周期管理
│   │   ├── gc.go               #     空闲房间垃圾回收
│   │   └── recovery.go         #     房间崩溃恢复
│   ├── collab/                  #   协同编辑核心算法
│   │   ├── operation.go        #     协同操作定义
│   │   ├── operation_type.go   #     操作类型枚举
│   │   ├── applier.go          #     操作应用器（将操作应用到画布）
│   │   ├── validator.go        #     操作合法性校验
│   │   ├── conflict.go         #     冲突检测与解决
│   │   ├── dedup.go            #     操作去重（幂等性保证）
│   │   ├── versioned_field.go  #     带版本号的字段（版本管理）
│   │   ├── tombstone.go        #     墓碑机制（处理删除冲突）
│   │   └── undo.go             #     撤回/重做支持
│   ├── snapshot/                #   快照管理（画布数据持久化）
│   │   ├── snapshot.go         #     快照数据结构
│   │   ├── manager.go          #     快照管理器（创建/恢复/清理）
│   │   ├── codec.go            #     快照编解码（序列化/反序列化）
│   │   ├── manifest.go         #     快照清单（元数据管理）
│   │   └── minio_store.go      #     基于 MinIO 的快照存储
│   ├── wal/                     #   预写日志 (Write-Ahead Log)
│   │   ├── wal.go              #     WAL 核心接口定义
│   │   ├── record.go           #     日志记录结构
│   │   ├── codec.go            #     日志编解码
│   │   ├── file_wal.go         #     基于文件的 WAL 实现
│   │   ├── segment.go          #     日志段管理（分段存储/滚动）
│   │   ├── recovery.go         #     WAL 恢复（崩溃后数据恢复）
│   │   ├── checksum.go         #     校验和（数据完整性验证）
│   │   └── rust_client.go      #     Rust 高性能 WAL 客户端绑定
│   ├── protocol/                #   通信协议编解码
│   │   ├── codec.go            #     编解码器接口
│   │   ├── json_codec.go       #     JSON 编解码实现
│   │   ├── proto_codec.go      #     Protobuf 编解码实现
│   │   └── message.go          #     通信消息结构体
│   ├── auth/                    #   认证与授权
│   │   ├── jwt.go              #     JWT 令牌生成与验证
│   │   ├── password.go         #     密码哈希与校验
│   │   ├── permission.go       #     权限检查（角色/资源）
│   │   └── service.go          #     认证服务（登录/注册/刷新）
│   ├── repository/              #   数据访问层（数据库操作）
│   │   ├── db.go               #     数据库连接初始化
│   │   ├── user_repo.go        #     用户数据访问
│   │   ├── room_repo.go        #     房间数据访问
│   │   ├── member_repo.go      #     成员数据访问
│   │   ├── snapshot_repo.go    #     快照元数据访问
│   │   └── file_repo.go        #     文件数据访问
│   ├── model/                   #   数据模型定义
│   │   ├── user.go             #     用户模型
│   │   ├── room.go             #     房间模型
│   │   ├── member.go           #     房间成员模型
│   │   ├── file.go             #     文件模型
│   │   └── snapshot.go         #     快照模型
│   ├── service/                 #   业务逻辑层
│   │   ├── auth_service.go     #     认证业务逻辑
│   │   ├── room_service.go     #     房间业务逻辑
│   │   ├── file_service.go     #     文件业务逻辑
│   │   └── user_service.go     #     用户业务逻辑
│   ├── storage/                 #   对象存储抽象
│   │   ├── minio.go            #     MinIO 客户端初始化
│   │   └── object_store.go     #     对象存储接口定义
│   ├── redis/                   #   Redis 客户端
│   │   ├── client.go           #     Redis 连接管理
│   │   ├── limiter.go          #     限流器（基于 Redis）
│   │   └── presence.go         #     在线状态管理
│   ├── observability/           #   可观测性（监控/日志/追踪）
│   │   ├── logger.go           #     日志记录器
│   │   ├── metrics.go          #     指标采集（Prometheus）
│   │   ├── tracing.go          #     分布式追踪（OpenTelemetry）
│   │   └── middleware.go       #     可观测性中间件
│   └── errors/                  #   错误定义与处理
│       ├── code.go             #     错误码定义
│       └── errors.go           #     错误类型与包装
├── migrations/                  # # 数据库迁移脚本
│   ├── 001_init_users.sql      #   用户表初始化
│   ├── 002_init_rooms.sql      #   房间表初始化
│   ├── 003_init_snapshots.sql  #   快照表初始化
│   └── 004_init_files.sql      #   文件表初始化
├── pkg/                         # # 公共工具包
│   ├── uuid/                    #   UUID 生成器
│   ├── clock/                   #   时钟抽象（便于测试时 mock）
│   └── retry/                   #   重试机制（指数退避）
└── tests/                       # # 测试目录
    ├── integration/             #   集成测试
    ├── room/                    #   房间模块测试
    ├── wal/                     #   WAL 模块测试
    ├── collab/                  #   协同编辑测试
    └── websocket/               #   WebSocket 测试