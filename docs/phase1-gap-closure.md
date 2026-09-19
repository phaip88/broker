# Phase 0/1 缺口关闭基线

远端 Workflow Run #35443967712 已证明 v1 基础闭环可运行。本文件冻结下一阶段的实现边界。

## 已落地的协议源

`proto/broker/v1/broker.proto` 是统一协议的唯一源文件。REST 保留为兼容入口，但新能力必须先进入 protobuf 模型，再提供 REST/Connect 映射。生成代码不得手工编辑；CI 应在安装 `buf` 后执行 lint、breaking 检查和代码生成差异检查。

## 领域模型

- `Project` 是隔离和权限边界。
- `OriginGraph` 记录主站、认证、API、CDN、支付等节点和允许的边。
- `Lease` 绑定 project、graph revision、执行模式、预算和生命周期。
- `Event` 使用 `trace_id`、`request_ref`、`action_ref` 关联浏览器、代理和编排动作。

旧 `Target` REST 模型只作为兼容层，不再作为新功能的主模型。

## 出网与 Relay 实现顺序

1. `BrowserRelay.Open` 为每个 lease 建立长连接。
2. relay 只接受来自 Browser Worker 的 CONNECT/HTTP 流，不允许浏览器绕过 relay 直连。
3. relay 在发送首个上游字节前校验 lease、graph revision、origin、方法、预算和执行模式。
4. headers/body 采用流式 frame；默认不把原始 body 写入日志。
5. 取消 lease 时向所有关联 stream 发送 RESET，并关闭上游连接。

在 relay 完成前，Browser Controller 只能作为开发模式适配器，不得声称实现了“物理阻断直连”。

## 流关联

Browser Worker 生成 `trace_id`；每个 Playwright request/response 和 mitmproxy flow 映射到同一 `request_ref`。CDP 事件用 `action_ref` 关联页面动作。Broker 只持久化脱敏事件，原始 body 留在短生命周期执行器内。

## 会话交接与取消

人工登录/MFA 使用显式状态：`WAITING_FOR_USER`、`READY`、`RUNNING`、`CANCEL_REQUESTED`、`CANCELLED`、`COMPLETED`。MFA 仅允许人工完成，Broker 不自动绕过验证。多 Tab/Popup 每个 target 具有独立 ref，但共享 lease 和 trace context。取消必须从 Broker 传播到 Browser Worker、Relay 和 mitmproxy adapter，并且幂等。

## 验收条件

- protobuf lint、生成和兼容性检查通过；
- 两个 origin 节点可在同一 lease 内按 graph edge 出网，未登记节点被拒绝；
- 浏览器关闭直连代理后页面仍能通过 relay 完成导航和 XHR；
- Playwright、CDP、mitmproxy 对同一请求得到同一 `request_ref`；
- 人工登录交接、Popup、取消和重复取消均有可查询状态；
- 断开 Worker 后 lease 超时或取消，所有 relay stream 均被回收；
- CI 中执行 Go、Node、协议 lint、浏览器 smoke test 和 Docker compose smoke test。
