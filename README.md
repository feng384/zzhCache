# zzhCache
zzhCache - 一个基于Go的轻量级分布式缓存系统

![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?logo=go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

zzhCache是一个受memcache和groupcache启发的分布式缓存系统，旨在提供高效的内存缓存解决方案，支持分布式节点扩展和智能缓存管理。

✨ 特性

• 分布式节点：支持多节点部署，通过一致性哈希实现请求路由

• LRU淘汰策略：自动移除最近最少使用的缓存项

• 防止缓存击穿：单飞机制（Singleflight）确保并发请求只执行一次数据加载

• HTTP接口：提供RESTful API进行缓存操作

• Protobuf序列化：高效二进制协议通信

• 可扩展架构：易于添加新的节点和缓存策略


🚀 快速开始

前置要求
• Go 1.19+

• Git


安装 & 运行
```bash
# 克隆仓库
git clone https://github.com/yourusername/zzhcache.git
cd zzhcache

# 构建并启动集群（3节点+API服务）
chmod +x run.sh
./run.sh
```

测试查询
```bash
curl "http://localhost:9999/api?key=Tom"
# 预期返回: 630
```

📖 使用指南

启动参数
| 参数    | 描述                  | 默认值  |
|---------|----------------------|---------|
| -port   | 服务端口号            | 8001    |
| -api    | 是否启用API服务       | false   |

示例：
```bash
# 启动缓存节点（端口8002）
go run main.go -port=8002

# 启动带API服务的节点
go run main.go -port=8003 -api=true
```

API接口
GET /api
• 参数: `key` - 要查询的缓存键

• 示例: `http://localhost:9999/api?key=Jack`


响应格式：
```json
"589"
```

🛠 部署架构

![架构图](https://via.placeholder.com/800x400.png?text=zzhCache+Architecture)

1. 客户端通过API Gateway访问
2. 一致性哈希环路由请求到对应节点
3. 节点优先检查LRU缓存
4. 未命中时通过数据加载器获取数据
5. 支持横向扩展缓存节点

📂 项目结构
```
zzhcache/
├── main.go                 # 服务入口
├── zzhcache/               # 核心实现
│   ├── lru/               # LRU缓存实现
│   ├── consistenthash/    # 一致性哈希算法
│   ├── singleflight/      # 单飞机制
│   └── zzhcachepb/        # Protobuf定义
├── run.sh                  # 集群启动脚本
└── README.md               # 项目文档
```

🔧 性能优化

• 缓存预热：启动时加载高频数据

• 批量加载：支持多键查询优化

• TTL支持：计划中的过期时间功能


🤝 如何贡献
欢迎提交Issue和PR！请确保：
1. 通过所有单元测试
2. 更新相关文档
3. 保持代码风格一致

```bash
# 运行测试套件
go test -v ./...
```

📄 协议
本项目基于 [MIT License](LICENSE) 开源。

---

*提示：zzhCache当前为v0.1预览版，生产环境使用建议进行充分测试*