# go-outboxrelay

Go 实现的事务发件箱中继：业务先 `Append` 事件到本地 outbox，再由 `RelayOnce` / `RelayContext` 经 HTTP 投递到下游。状态机 pending→sending→done/dead，支持退避重试与管理页。

## 运行

```bash
go test ./... -count=1
go run ./cmd/outboxd -addr :8091 -web web -target http://127.0.0.1:9999/hook
```

浏览器打开 `http://localhost:8091/`：追加事件、触发中继、查看 pending 与统计。

## 库面

```go
o := outboxrelay.New(outboxrelay.WithTargetURL("https://example.com/hook"))
id, err := o.Append("order.created", []byte(`{"id":1}`), nil)
results, err := o.RelayContext(ctx, 10)
_ = o.ListPending(20)
_ = o.Close()
```
