基于 Go 实现的冷库氨制冷机组群控平台项目，一款后端服务，完成库房氨制冷机组的加载、供液、化霜、库温速率窗口、氨泄漏联锁与冷负荷调度控制。

## 运行

本项目使用 Go 1.23，依赖已随 vendor 目录离线携带。

```bash
go build -mod=vendor ./cmd/coldstore
COLDSTORE_DATA=./coldstore-data.json COLDSTORE_ADDR=:8080 ./coldstore
```

服务启动后提供 JSON API 控制台，健康检查地址为 `/health`，运行状态地址为 `/status`。
