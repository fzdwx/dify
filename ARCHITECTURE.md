# Dify Go Client 架构说明

## 认证架构

Dify Go 客户端使用双重认证机制来处理不同类型的 API 调用：

```
┌─────────────────────────────────────────────────────────────┐
│                    Dify Go Client                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐              ┌─────────────────┐      │
│  │  Console Client │              │ Datasets Client │      │
│  │                 │              │                 │      │
│  │ • Access Token  │              │ • API Key       │      │
│  │ • Refresh Token │              │ • No Expiry     │      │
│  │ • Auto Refresh  │              │ • No Refresh    │      │
│  └─────────────────┘              └─────────────────┘      │
│           │                                 │               │
│           │                                 │               │
│           ▼                                 ▼               │
│  ┌─────────────────┐              ┌─────────────────┐      │
│  │  Console APIs   │              │  Datasets APIs  │      │
│  │                 │              │                 │      │
│  │ /console/api/*  │              │ /v1/datasets/*  │      │
│  │                 │              │                 │      │
│  │ • Login         │              │ • Create Dataset│      │
│  │ • Get API Keys  │              │ • Upload Files  │      │
│  │ • Create Keys   │              │ • Manage Docs   │      │
│  └─────────────────┘              └─────────────────┘      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 认证流程

### 1. 初始化阶段

```go
client, err := dify.NewClient(baseUrl, email, password)
```

1. **用户登录** (`/console/api/login`)
   - 发送 email + password
   - 获取 access_token + refresh_token

2. **获取 Datasets API Key** (`/console/api/datasets/api-keys`)
   - 使用 access_token 认证
   - 获取现有 API key 或创建新的

3. **创建双客户端**
   - Console Client: 使用 access_token (会过期)
   - Datasets Client: 使用 datasets API key (不过期)

### 2. 运行时阶段

#### Console API 调用
```go
// 内部调用，用户不直接使用
client.RefreshDatasetAPIKey() // 调用 /console/api/datasets/api-keys
```

- 使用 `consoleClient` 
- 如果遇到 401 错误：
  1. 自动调用 `/console/api/refresh-token`
  2. 更新 access_token 和 refresh_token
  3. 重试原始请求

#### Datasets API 调用
```go
// 用户直接调用
client.CreateEmptyDataset(ctx, req)  // 调用 /v1/datasets
client.CreateByFile(ctx, req)        // 调用 /v1/datasets/{id}/document/create-by-file
```

- 使用 `datasetsClient`
- 使用 datasets API key 认证
- **不需要** refresh token 机制
- API key 通常不会过期

## 错误处理

### Console API 错误处理
```go
func (c *client) executeConsoleWithRetry(requestFunc func() (*resty.Response, error)) (*resty.Response, error) {
    // 第一次尝试
    response, err := requestFunc()
    
    // 检查 401 错误
    if response.StatusCode() == 401 {
        // 刷新 access token
        c.refreshAccessToken()
        // 重试请求
        response, err = requestFunc()
    }
    
    return response, err
}
```

### Datasets API 错误处理
```go
// 直接返回错误，不进行 token 刷新
func (c *client) CreateEmptyDataset(ctx context.Context, req *CreateEmptyDatasetRequest) (*Response[CreateEmptyDatasetResponse], error) {
    response, err := c.datasets().Post("/datasets")
    return buildResponse(response, resp), err
}
```

## 关键设计原则

1. **职责分离**: Console API 和 Datasets API 使用不同的认证机制
2. **自动化**: 用户无需手动管理 token 刷新
3. **透明性**: refresh token 机制对用户完全透明
4. **效率**: 只对需要的 API (Console API) 应用 refresh 机制
5. **健壮性**: 完整的错误处理和重试机制

## 使用建议

- **正常使用**: 直接调用 datasets 相关方法，无需关心认证细节
- **长期运行**: 客户端会自动处理 token 过期，适合长期运行的应用
- **错误处理**: 关注业务逻辑错误，认证错误会自动处理
- **手动刷新**: 通常不需要，但可以调用 `RefreshDatasetAPIKey()` 强制刷新
