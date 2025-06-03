# Dify Golang Client

dify 的 Golang 客户端，用于管理数据集、创建应用等。

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

## 许可证

MIT License