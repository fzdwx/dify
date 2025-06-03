# Dify Golang Client

一个功能完整的 Dify API Go 客户端，支持用户名密码登录、自动 token 刷新和 datasets 管理。

## 特性

- ✅ 用户名密码登录认证
- ✅ 自动获取/创建 datasets API key
- ✅ 自动 refresh token 机制
- ✅ 支持创建空数据集
- ✅ 支持通过文件创建文档
- ✅ 完整的错误处理和重试机制

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/fzdwx/dify"
)

func main() {
    // 使用用户名和密码创建客户端
    // 客户端会自动登录并获取或创建 datasets API key
    client, err := dify.NewClient("http://your-dify-host", "your-email", "your-password")
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }

    ctx := context.Background()

    // 创建数据集
    uniqueName := fmt.Sprintf("测试数据集_%d", time.Now().Unix())
    resp, err := client.CreateEmptyDataset(ctx, &dify.CreateEmptyDatasetRequest{
        Name:              uniqueName,
        Description:       "测试数据集",
        IndexingTechnique: dify.IndexingTechniqueEconomy,
        Permission:        dify.DatasetPermissionAllTeamMembers,
        Provider:          dify.DatasetProviderVendor,
        RetrievalModel: dify.RetrievalModel{
            SearchMethod:    dify.RetrievalModelSearchMethodHybridSearch,
            RerankingEnable: true,
            TopK:            10,
        },
    })

    if err != nil {
        log.Fatal("创建数据集失败:", err)
    }

    fmt.Printf("成功创建数据集: %s (ID: %s)\n", resp.Result.Name, resp.Result.ID)
}
```

## 认证机制

### 自动登录和 Token 管理

客户端使用以下认证流程：

1. **用户登录**: 使用 email 和 password 登录获取 access token 和 refresh token
2. **获取 Datasets API Key**: 自动获取现有的或创建新的 datasets API key
3. **自动刷新**: 当 console access token 过期时自动使用 refresh token 刷新

### Token 过期处理

客户端使用两种不同的认证机制：

#### Console API (`/console/api/` 开头的接口)
- 使用 **access token** 认证
- Token 会过期，返回 `{"code": "unauthorized", "message": "Token has expired.", "status": 401}`
- **自动处理**: 遇到 401 错误时自动使用 refresh token 刷新并重试
- 涉及的操作：获取/创建 datasets API keys

#### Datasets API (`/v1/datasets/` 开头的接口)
- 使用 **datasets API key** 认证
- API key 通常不会过期
- **无需 refresh token**: 直接使用 datasets API key，不涉及 token 刷新
- 涉及的操作：创建数据集、上传文档等

## API 参考

### 创建客户端

```go
client, err := dify.NewClient(baseUrl, email, password string) (Client, error)
```

### 数据集操作

```go
// 创建空数据集
resp, err := client.CreateEmptyDataset(ctx, req)

// 通过文件创建文档
resp, err := client.CreateByFile(ctx, req)

// 手动刷新 datasets API key（可选）
err := client.RefreshDatasetAPIKey()
```

## 错误处理

客户端内置了完整的错误处理和重试机制：

- 自动检测 401 错误并刷新 token
- 详细的错误信息和状态码
- 网络错误的适当处理

## 许可证

MIT License