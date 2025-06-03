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
	// 同时支持自动 refresh token 功能
	client, err := dify.NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}

	ctx := context.Background()

	// 创建空的数据集
	uniqueName := fmt.Sprintf("示例数据集_%d", time.Now().Unix())
	resp, err := client.CreateEmptyDataset(ctx, &dify.CreateEmptyDatasetRequest{
		Name:              uniqueName,
		Description:       "这是一个示例数据集，支持自动 token 刷新",
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

	if !resp.IsSuccess() {
		log.Fatal("创建数据集失败:", resp.Message)
	}

	fmt.Printf("成功创建数据集: %s (ID: %s)\n", resp.Result.Name, resp.Result.ID)

	// 如果需要手动刷新 dataset API key（通常不需要，因为有自动重试机制）
	// err = client.RefreshDatasetAPIKey()
	// if err != nil {
	//     log.Fatal("刷新 dataset API key 失败:", err)
	// }
	// fmt.Println("成功刷新 dataset API key")
}
