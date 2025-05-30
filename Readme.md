# app.yaml 配置文件

```yaml
App:
  port: 7111
Alter:
  Grafana:
    hook_uri: /hook/grafana  # 用于写入Grafana配置
  AlertManager:
    hook_uri: /hook/alertManager  # 用于写入AlertManager配置
Notifier:
  DingTalk:
    hook_uri: /dingTalk  # 用于直接发送通知的http接口
    groups:
      - name: 开发告警组
        token: a token
        secret: a secret
  ShowDoc:
    hook_uri: /showDoc  # 用于直接发送通知的http接口
    tokens:
      - name: Sven
        token: foo
      - name: 老张
        token: bar

```
