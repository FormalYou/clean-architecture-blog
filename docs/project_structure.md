# 项目结构

本文档详细描述了本项目的目录和文件结构，旨在帮助开发人员快速理解代码组织方式、各个组件的职责以及它们之间的相互关系。

```
.
├── .golangci.yml
├── Makefile
├── README.md
├── api/
│   └── openapi.yaml
├── cmd/
│   └── server/
│       ├── main.go
│       └── option/
│           └── option.go
├── configs/
│   └── config.yaml
├── docs/
│   ├── SLO.md
│   └── project_structure.md
├── domain/
│   ├── article.go
│   ├── audit_event.go
│   ├── comment.go
│   ├── tag.go
│   └── user.go
├── go.mod
├── go.sum
├── internal/
│   ├── application/
│   │   ├── contracts/
│   │   ├── repository/
│   │   └── usecase/
│   ├── errorx/
│   ├── infrastructure/
│   │   ├── auth/
│   │   ├── cache/
│   │   ├── config/
│   │   ├── log/
│   │   └── persistence/
│   └── interfaces/
│       └── http/
│           ├── dto/
│           ├── handler/
│           └── middleware/
├── scripts/
│   ├── test-e2e.sh
│   ├── test-integration.sh
│   └── test-unit.sh
└── tests/
    └── e2e/
```

### 根目录

*   `.golangci.yml`: [golangci-lint](https://golangci-lint.run/) 的配置文件。它定义了代码风格和质量检查的规则，用于确保代码库的一致性和健壮性。
*   `Makefile`: 包含了一系列用于自动化常见开发任务的命令，例如构建、测试、运行和部署应用程序。
*   `README.md`: 项目的入口文档，提供了项目概述、安装指南、使用方法和贡献准则等基本信息。
*   `go.mod`: Go 模块文件，定义了项目的模块路径和依赖项。
*   `go.sum`: 记录了特定依赖版本的校验和，以确保构建的可重现性和安全性。

### `api/`

此目录存放 API 契约和文档。

*   `openapi.yaml`: OpenAPI (Swagger) 规范文件，用以定义 RESTful API 的端点、请求/响应格式和数据模型。

### `cmd/`

此目录包含项目的主要应用程序入口。

*   `server/`: 存放 Web 服务器相关代码。
    *   `main.go`: 应用程序的主入口点。负责初始化配置、日志、数据库连接、依赖注入和启动 HTTP 服务器。
    *   `option/`: 存放服务器启动选项和配置。
        *   `option.go`: 定义了用于配置服务器启动的结构体和函数。

### `configs/`

存放应用程序的配置文件。

*   `config.yaml`: 默认的配置文件，包含了数据库连接信息、服务器端口、日志级别等配置。

### `docs/`

包含项目的相关文档。

*   `SLO.md`: 服务级别目标 (Service Level Objectives) 文档，定义了服务的可用性和性能目标。
*   `project_structure.md`: (本文) 详细描述了项目的文件和目录结构。

### `domain/`

包含核心业务逻辑和实体，这是领域驱动设计 (DDD) 中的领域层。此处的代码不应依赖任何外部框架或实现细节。

*   `article.go`: 定义了文章 (Article) 的领域模型/实体。
*   `audit_event.go`: 定义了审计事件 (Audit Event) 的领域模型。
*   `comment.go`: 定义了评论 (Comment) 的领域模型。
*   `tag.go`: 定义了标签 (Tag) 的领域模型。
*   `user.go`: 定义了用户 (User) 的领域模型。

### `internal/`

包含所有不对外暴露的私有应用程序代码。Go 语言的这个特性可以防止其他项目导入此目录下的包。

*   **`application/`**: 应用层，编排领域对象执行业务用例。
    *   `contracts/`: 定义了应用层与基础设施层之间的接口（契约），例如日志、认证和审计服务。
    *   `repository/`: 定义了仓储接口，用于抽象数据持久化逻辑。
    *   `usecase/`: 包含了具体的业务用例（或称交互器），实现了应用的核心功能。
*   **`errorx/`**: 包含自定义的错误类型和错误处理帮助函数。
*   **`infrastructure/`**: 基础设施层，提供了应用层所需服务的具体实现，例如数据库、缓存、认证等。
    *   `auth/`: 包含了认证和授权的具体实现（例如 JWT）。
    *   `cache/`: 提供了缓存服务的实现（例如 Redis）。
    *   `config/`: 负责加载和解析配置文件（例如 Viper）。
    *   `log/`: 提供了日志服务的具体实现（例如 Zap）。
    *   `persistence/`: 实现了数据持久化逻辑，通常是对仓储接口的具体实现（例如 GORM）。
*   **`interfaces/`**: 接口层（也称为表示层），负责与外部系统进行交互。
    *   `http/`: 包含了 HTTP 服务相关代码。
        *   `dto/`: 数据传输对象 (Data Transfer Objects)，用于在接口层和应用层之间传输数据。
        *   `handler/`: HTTP 处理器，负责解析请求、调用应用层用例并返回响应。
        *   `middleware/`: HTTP 中间件，用于处理横切关注点，如认证、日志、错误恢复等。

### `scripts/`

存放用于支持开发、测试和部署流程的脚本。

*   `test-e2e.sh`: 端到端测试脚本。
*   `test-integration.sh`: 集成测试脚本。
*   `test-unit.sh`: 单元测试脚本。

### `tests/`

包含自动化测试代码。

*   `e2e/`: 端到端 (End-to-End) 测试，用于模拟真实用户场景，测试整个应用程序的流程。
