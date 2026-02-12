# Go 仓库管理系统

这是一个将原 [JavaFX 仓库管理系统](https://github.com/wzomg/jfoenix-WarehouseMS)重构为 Go Web 应用的版本，提供登录/注册、商品增删改查与库存统计功能，并包含一键启动与停止脚本。

## 运行环境
- Go 1.23+
- Docker 24+

## 配置
编辑 config/config.yaml，设置数据库连接与端口。
默认数据库容器映射端口为 3307。

## 启动与停止
```bash
./scripts/control.sh start
./scripts/control.sh status
./scripts/control.sh stop
```

## 访问地址
- http://localhost:8081

## 技术栈
- Go 1.25
- Gin（HTTP 服务）
- Gorm + MySQL（数据访问）
- Viper（配置管理）
- Zap（日志）
- Docker Compose（本地数据库容器）
- Bootstrap 5（前端样式）

## 注意事项
- 数据持久化：重复启动不会丢数据，只有执行 `docker compose down -v` 才会清空数据库卷
- 初始化数据：首次启动会加载 `db/init.sql` 中的示例数据
- 数据库连接：端口在 `config/config.yaml` 中配置，默认容器映射为 3307
- 撤销接口：最近一次库存变更可通过 POST `/api/goods/undo` 撤销

## 设计模式
- 单例模式：数据库连接与日志实例
  - 目的：确保全局只有一份连接与日志器，避免重复初始化与资源浪费
  - 位置：
    - 数据库单例：[db.go](internal/repository/db.go#L36-L47)
    - 日志单例：[logger.go](internal/infra/logger.go#L9-L20)
  - 使用方式：
    - 统一通过 GetDB / GetLogger 获取实例，再注入到仓储、服务与路由
    - 组装入口：[main.go](cmd/server/main.go#L47-L75)

- 原型模式：商品对象克隆
  - 目的：批量新增时，基于已有对象快速复制，避免散落的手工复制逻辑
  - 位置：
    - 原型方法：[goods.go](internal/model/goods.go#L11-L14)
  - 使用方式：
    - 新增商品时对传入对象 Clone，再入库：[goods_service.go](internal/service/goods_service.go#L45-L59)

- 代理模式：商品仓储代理
  - 目的：在不修改真实仓储的前提下，追加日志、监控或缓存等横切逻辑
  - 位置：
    - 代理实现：[goods_proxy.go](internal/repository/goods_proxy.go#L9-L52)
    - 真实仓储与接口：[goods_repository.go](internal/repository/goods_repository.go#L9-L49)、[goods_store.go](internal/repository/goods_store.go#L1-L13)
  - 使用方式：
    - main 中注入代理代替真实仓储：[main.go](cmd/server/main.go#L63-L71)

- 观察者模式：库存变更事件通知
  - 目的：解耦“库存变更”与“审计/通知/统计”等副作用
  - 位置：
    - 事件总线与观察者：[events.go](internal/service/events.go#L12-L73)
  - 使用方式：
    - 变更后发布事件：[goods_service.go](internal/service/goods_service.go#L45-L105)
    - 注册观察者：[main.go](cmd/server/main.go#L65-L68)

- 备忘录模式：库存变更撤销
  - 目的：为“出库/删除”保留快照，实现撤销与回滚
  - 位置：
    - 备忘录与看护者：[memento.go](internal/service/memento.go#L10-L53)
  - 使用方式：
    - 变更前保存快照，变更后可撤销：[goods_service.go](internal/service/goods_service.go#L62-L105)
    - 撤销接口：POST `/api/goods/undo`：[goods_handler.go](internal/api/goods_handler.go#L79-L85)、[router.go](internal/api/router.go#L29-L36)
