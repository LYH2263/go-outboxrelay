// Package outboxrelay 实现事务发件箱：业务先 Append 事件到本地存储，
// 再由 RelayOnce / RelayContext 经 HTTP 投递到下游，状态 pending→sending→done/dead。
package outboxrelay
