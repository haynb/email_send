# 邮件发送服务

一个轻量级邮件发送HTTP服务，可以通过HTTP接口使用163邮箱发送电子邮件。

## 功能

- 通过HTTP API发送邮件
- 支持纯文本和HTML格式邮件
- 支持多收件人
- 使用Docker容器部署
- 轻量级设计
- API密钥身份验证保护

## 快速开始

### 配置

1. 编辑`.env`文件，填入您的163邮箱信息：

```
SMTP_SERVER=smtp.163.com
SMTP_PORT=465
SENDER_EMAIL=你的邮箱@163.com
SENDER_NAME=邮件服务
AUTH_PASSWORD=您的授权密码
PORT=8080
API_KEY=your-secure-api-key-here
```

请确保将 `API_KEY` 修改为一个安全的、难以猜测的值。

### 部署

使用Docker Compose启动服务：

```bash
docker-compose up -d
```

## API使用

### 身份验证

所有API请求（除了健康检查）都需要在HTTP头部包含有效的API密钥：

```
X-API-Key: your-secure-api-key-here
```

### 发送邮件

**请求**:

```
POST /send-email
Content-Type: application/json
X-API-Key: your-secure-api-key-here

{
  "to": ["收件人1@example.com", "收件人2@example.com"],
  "subject": "邮件主题",
  "body": "邮件内容",
  "is_html": false
}
```

如果需要发送HTML格式的邮件，将`is_html`设置为`true`，并在`body`中提供HTML内容。

**响应**:

```json
{
  "success": true,
  "message": "邮件发送成功"
}
```

### 健康检查

```
GET /health
```

健康检查接口不需要API密钥。

## 示例

使用curl发送邮件：

```bash
curl -X POST http://localhost:8080/send-email \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secure-api-key-here" \
  -d '{
    "to": ["收件人@example.com"],
    "subject": "服务器告警",
    "body": "CPU使用率超过90%！",
    "is_html": false
  }'
```

## 注意事项

- 请确保您的163邮箱已经开启POP3/SMTP服务
- 使用授权密码而不是登录密码
- 使用SMTP服务发送邮件可能会有每日发送数量限制
- 应当设置一个强API密钥来保护您的服务 